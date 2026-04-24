#!/usr/bin/env python3
"""Compare a tmp*.mcfunction source command with the auto-generated output.

Parses the first non-comment command line of each file (either `give ...` or
`summon ...`), decodes the SNBT payload, and verifies that every key/value
present in the tmp command is also present in the generated command. Extra
keys in the generated output (e.g., `DeathLootTable`, injected maf tags,
`minecraft:custom_data`) are tolerated — the script only flags drift that
would indicate savedata under-coverage.

Exit codes
    0 — every tmp leaf is present in generated (possibly with
        `minecraft:` prefix tolerance or text-component wrapping).
    1 — discrepancies found; they are listed on stdout.
    2 — parse error (the command could not be interpreted).
"""

from __future__ import annotations

import argparse
import re
import sys
from pathlib import Path
from typing import Any


class ParseError(Exception):
    pass


# --- SNBT parser -----------------------------------------------------------

_SCALAR_RE = re.compile(r"^(-?\d+(?:\.\d+)?(?:[eE][+-]?\d+)?)([bBsSlLfFdD])?$")
_KEY_TERMINATORS = ':,{}[]" \t\n\r'
_SCALAR_TERMINATORS = ',{}[]" \t\n\r'


class SNBTParser:
    """Forgiving SNBT → Python value parser for the subset we see in MC."""

    def __init__(self, text: str):
        self.text = text
        self.pos = 0

    def parse(self) -> Any:
        self._skip_ws()
        val = self._parse_value()
        self._skip_ws()
        return val

    def _eof(self) -> bool:
        return self.pos >= len(self.text)

    def _peek(self, offset: int = 0) -> str:
        idx = self.pos + offset
        return self.text[idx] if idx < len(self.text) else ""

    def _skip_ws(self) -> None:
        while not self._eof() and self.text[self.pos].isspace():
            self.pos += 1

    def _parse_value(self) -> Any:
        self._skip_ws()
        c = self._peek()
        if c == "{":
            return self._parse_object()
        if c == "[":
            return self._parse_array()
        if c in ('"', "'"):
            return self._parse_quoted_string(c)
        return self._parse_scalar()

    def _parse_object(self) -> dict:
        self.pos += 1  # consume '{'
        obj: dict = {}
        self._skip_ws()
        if self._peek() == "}":
            self.pos += 1
            return obj
        while True:
            self._skip_ws()
            key = self._parse_key()
            self._skip_ws()
            if self._peek() != ":":
                raise ParseError(f"expected ':' after key {key!r} at pos {self.pos}")
            self.pos += 1
            obj[key] = self._parse_value()
            self._skip_ws()
            if self._peek() == ",":
                self.pos += 1
                self._skip_ws()
                if self._peek() == "}":  # tolerate trailing comma
                    self.pos += 1
                    return obj
                continue
            if self._peek() == "}":
                self.pos += 1
                return obj
            raise ParseError(f"expected ',' or '}}' at pos {self.pos}")

    def _parse_array(self) -> list:
        self.pos += 1  # consume '['
        self._skip_ws()
        # Typed arrays: [B; ...] / [I; ...] / [L; ...]
        if (
            self.pos + 1 < len(self.text)
            and self.text[self.pos] in "BIL"
            and self.text[self.pos + 1] == ";"
        ):
            self.pos += 2
            self._skip_ws()
        arr: list = []
        if self._peek() == "]":
            self.pos += 1
            return arr
        while True:
            arr.append(self._parse_value())
            self._skip_ws()
            if self._peek() == ",":
                self.pos += 1
                continue
            if self._peek() == "]":
                self.pos += 1
                return arr
            raise ParseError(f"expected ',' or ']' at pos {self.pos}")

    def _parse_quoted_string(self, quote: str) -> str:
        self.pos += 1
        chars: list[str] = []
        while not self._eof():
            c = self.text[self.pos]
            if c == "\\":
                self.pos += 1
                if self._eof():
                    break
                nxt = self.text[self.pos]
                chars.append({"n": "\n", "t": "\t", "r": "\r"}.get(nxt, nxt))
                self.pos += 1
                continue
            if c == quote:
                self.pos += 1
                return "".join(chars)
            chars.append(c)
            self.pos += 1
        raise ParseError("unterminated string literal")

    def _parse_key(self) -> str:
        c = self._peek()
        if c in ('"', "'"):
            return self._parse_quoted_string(c)
        start = self.pos
        while not self._eof() and self.text[self.pos] not in _KEY_TERMINATORS:
            self.pos += 1
        if self.pos == start:
            raise ParseError(f"empty key at pos {self.pos}")
        return self.text[start:self.pos]

    def _parse_scalar(self) -> Any:
        start = self.pos
        while not self._eof() and self.text[self.pos] not in _SCALAR_TERMINATORS:
            self.pos += 1
        tok = self.text[start:self.pos]
        if not tok:
            raise ParseError(f"empty scalar at pos {self.pos}")
        if tok in ("true", "True"):
            return True
        if tok in ("false", "False"):
            return False
        m = _SCALAR_RE.match(tok)
        if m:
            num_text = m.group(1)
            suffix = m.group(2)
            force_float = suffix in ("f", "F", "d", "D") or "." in num_text or "e" in num_text.lower()
            try:
                return float(num_text) if force_float else int(num_text)
            except ValueError:
                pass
        return tok


# --- components parser (give) ---------------------------------------------


def parse_components(text: str) -> dict:
    """Parse `key=val,key=val` at the top level, where val is SNBT."""
    i = 0
    n = len(text)
    out: dict = {}
    while i < n:
        while i < n and text[i] in " \t\n\r,":
            i += 1
        if i >= n:
            break
        key_start = i
        while i < n and text[i] != "=":
            i += 1
        key = text[key_start:i].strip()
        if not key:
            break
        if i >= n:
            raise ParseError(f"components: no '=' for key {key!r}")
        i += 1  # consume '='
        val_start = i
        depth_brace = 0
        depth_brack = 0
        in_str: str | None = None
        while i < n:
            c = text[i]
            if in_str:
                if c == "\\":
                    i += 2
                    continue
                if c == in_str:
                    in_str = None
                i += 1
                continue
            if c in ('"', "'"):
                in_str = c
                i += 1
                continue
            if c == "{":
                depth_brace += 1
            elif c == "}":
                depth_brace -= 1
            elif c == "[":
                depth_brack += 1
            elif c == "]":
                depth_brack -= 1
            elif c == "," and depth_brace == 0 and depth_brack == 0:
                break
            i += 1
        val_text = text[val_start:i].strip()
        out[key] = SNBTParser(val_text).parse()
    return out


# --- command extraction ----------------------------------------------------

_COMMENT_RE = re.compile(r"^\s*#")


def read_first_command_line(path: Path) -> str:
    for raw in path.read_text(encoding="utf-8").splitlines():
        s = raw.strip()
        if not s or _COMMENT_RE.match(s):
            continue
        return s
    raise ParseError(f"no command line in {path}")


def _normalize_id(raw: str) -> str:
    raw = raw.strip()
    return raw if ":" in raw else f"minecraft:{raw}"


def extract_payload(line: str) -> tuple[str, str, Any]:
    """Return (verb, id, payload) where payload is a dict for give/summon."""
    s = line.lstrip("/")
    parts = s.split(None, 1)
    if not parts:
        raise ParseError("empty command")
    verb = parts[0]
    rest = parts[1] if len(parts) > 1 else ""
    if verb == "give":
        selector_split = rest.split(None, 1)
        if len(selector_split) < 2:
            raise ParseError("give missing item")
        body = selector_split[1]
        bracket_idx = body.find("[")
        space_idx = body.find(" ")
        if bracket_idx == -1 or (space_idx != -1 and space_idx < bracket_idx):
            return verb, _normalize_id(body.split()[0]), {}
        item_id = body[:bracket_idx]
        i = bracket_idx
        depth = 0
        in_str: str | None = None
        while i < len(body):
            c = body[i]
            if in_str:
                if c == "\\":
                    i += 2
                    continue
                if c == in_str:
                    in_str = None
                i += 1
                continue
            if c in ('"', "'"):
                in_str = c
                i += 1
                continue
            if c == "[":
                depth += 1
            elif c == "]":
                depth -= 1
                if depth == 0:
                    break
            i += 1
        components_text = body[bracket_idx + 1 : i]
        return verb, _normalize_id(item_id), parse_components(components_text)
    if verb == "summon":
        tokens = rest.split(None, 4)
        if len(tokens) < 4:
            raise ParseError("summon missing coordinates")
        mob_type = tokens[0]
        nbt_text = tokens[4] if len(tokens) >= 5 else ""
        nbt = SNBTParser(nbt_text).parse() if nbt_text else {}
        return verb, _normalize_id(mob_type), nbt
    raise ParseError(f"unsupported verb {verb!r}")


# --- diffing ---------------------------------------------------------------


def _strip_prefix(k: str) -> str:
    return k.split(":", 1)[1] if k.startswith("minecraft:") else k


def _norm_primitive(v: Any) -> Any:
    if isinstance(v, bool):
        return int(v)
    if isinstance(v, float) and v.is_integer():
        return int(v)
    return v


def _leaves_equal(lv: Any, rv: Any) -> bool:
    if _norm_primitive(lv) == _norm_primitive(rv):
        return True
    try:
        if isinstance(lv, str) and isinstance(rv, (int, float)):
            return float(lv) == float(rv)
        if isinstance(rv, str) and isinstance(lv, (int, float)):
            return float(rv) == float(lv)
    except ValueError:
        pass
    return False


def _find_best_match(lv: Any, right_list: list) -> tuple[int | None, list]:
    best_idx: int | None = None
    best_diffs: list = []
    best_score: int | None = None
    for idx, rv in enumerate(right_list):
        diffs = diff_contains(lv, rv, "")
        if not diffs:
            return idx, []
        if best_score is None or len(diffs) < best_score:
            best_score = len(diffs)
            best_idx = idx
            best_diffs = diffs
    return best_idx, best_diffs


def diff_contains(left: Any, right: Any, path: str) -> list[tuple[str, str]]:
    """Return a list of (path, description) for every leaf in `left` not found in `right`."""
    diffs: list[tuple[str, str]] = []
    if isinstance(left, dict):
        if not isinstance(right, dict):
            diffs.append((path or ".", f"expected dict in generated, got {type(right).__name__}"))
            return diffs
        right_by_stripped = {_strip_prefix(k): k for k in right.keys()}
        for k, v in left.items():
            kn = _strip_prefix(k)
            target_key = k if k in right else right_by_stripped.get(kn)
            if target_key is None:
                diffs.append((f"{path}.{k}", f"missing in generated (tmp value = {v!r})"))
                continue
            diffs.extend(diff_contains(v, right[target_key], f"{path}.{k}"))
        return diffs
    if isinstance(left, list):
        if not isinstance(right, list):
            diffs.append((path or ".", f"expected list in generated, got {type(right).__name__}"))
            return diffs
        for i, lv in enumerate(left):
            idx, sub = _find_best_match(lv, right)
            if idx is None:
                diffs.append((f"{path}[{i}]", f"no matching element (tmp value = {lv!r})"))
            elif sub:
                diffs.append(
                    (
                        f"{path}[{i}]",
                        f"closest match in generated[{idx}] still has {len(sub)} diff(s); e.g. {sub[0]}",
                    )
                )
        return diffs
    # primitives
    if _leaves_equal(left, right):
        return diffs
    # text-component wrapping tolerance
    if isinstance(left, str) and isinstance(right, dict) and right.get("text") == left:
        return diffs
    if isinstance(right, str) and isinstance(left, dict) and left.get("text") == right:
        return diffs
    diffs.append((path or ".", f"value mismatch (tmp={left!r}, generated={right!r})"))
    return diffs


# --- entry point -----------------------------------------------------------


def main(argv: list[str] | None = None) -> int:
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("tmp", type=Path, help="Path to tmp*.mcfunction (input source)")
    ap.add_argument("generated", type=Path, help="Path to auto-generated .mcfunction")
    ap.add_argument(
        "--dump",
        action="store_true",
        help="Dump normalized tmp/generated payloads as JSON for eyeballing",
    )
    args = ap.parse_args(argv)

    try:
        tmp_line = read_first_command_line(args.tmp)
        gen_line = read_first_command_line(args.generated)
        tverb, tid, tpay = extract_payload(tmp_line)
        gverb, gid, gpay = extract_payload(gen_line)
    except ParseError as e:
        print(f"parse error: {e}", file=sys.stderr)
        return 2

    if args.dump:
        import json

        print("# tmp command")
        print(json.dumps({"verb": tverb, "id": tid, "payload": tpay}, ensure_ascii=False, indent=2))
        print()
        print("# generated command")
        print(json.dumps({"verb": gverb, "id": gid, "payload": gpay}, ensure_ascii=False, indent=2))
        print()

    rc = 0
    if tverb != gverb:
        print(f"verb mismatch: tmp={tverb!r} generated={gverb!r}")
        rc = 1
    if tid != gid:
        print(f"id mismatch: tmp={tid!r} generated={gid!r}")
        rc = 1

    diffs = diff_contains(tpay, gpay, "")
    if diffs:
        rc = 1
        print(f"Discrepancies ({len(diffs)}):")
        for p, desc in diffs:
            print(f"  - {p or '.'}: {desc}")

    if rc == 0:
        print(f"OK: every tmp key/value is present in generated ({args.tmp.name} ⊆ {args.generated.name})")
    return rc


if __name__ == "__main__":
    sys.exit(main())
