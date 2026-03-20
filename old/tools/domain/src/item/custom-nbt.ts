type TopLevelPair = {
	key: string;
	value: string;
};

const ENCHANTMENT_COMPOUND_PATTERN =
	/\bid\s*:\s*"([^"]+)"[\s\S]*?\blvl\s*:\s*([+-]?\d+)(?:[bBsSlL])?/;

function trimOuterBraces(value: string): string {
	const trimmed = value.trim();
	if (trimmed.startsWith("{") && trimmed.endsWith("}")) {
		return trimmed.slice(1, -1).trim();
	}
	return trimmed;
}

function splitTopLevelPairs(fragment: string): TopLevelPair[] {
	const src = trimOuterBraces(fragment);
	if (!src) {
		return [];
	}

	const pairs: TopLevelPair[] = [];
	let index = 0;

	while (index < src.length) {
		while (index < src.length && (src[index] === "," || /\s/.test(src[index]!))) {
			index += 1;
		}
		if (index >= src.length) {
			break;
		}

		const keyStart = index;
		let colonIndex = -1;
		let inSingle = false;
		let inDouble = false;
		let escaped = false;
		while (index < src.length) {
			const ch = src[index]!;
			if (escaped) {
				escaped = false;
				index += 1;
				continue;
			}
			if (ch === "\\") {
				escaped = true;
				index += 1;
				continue;
			}
			if (ch === "'" && !inDouble) {
				inSingle = !inSingle;
				index += 1;
				continue;
			}
			if (ch === '"' && !inSingle) {
				inDouble = !inDouble;
				index += 1;
				continue;
			}
			if (!inSingle && !inDouble && ch === ":") {
				colonIndex = index;
				break;
			}
			index += 1;
		}

		if (colonIndex < 0) {
			break;
		}

		const key = src.slice(keyStart, colonIndex).trim();
		index = colonIndex + 1;

		const valueStart = index;
		let braceDepth = 0;
		let bracketDepth = 0;
		inSingle = false;
		inDouble = false;
		escaped = false;

		while (index < src.length) {
			const ch = src[index]!;
			if (escaped) {
				escaped = false;
				index += 1;
				continue;
			}
			if (ch === "\\") {
				escaped = true;
				index += 1;
				continue;
			}
			if (ch === "'" && !inDouble) {
				inSingle = !inSingle;
				index += 1;
				continue;
			}
			if (ch === '"' && !inSingle) {
				inDouble = !inDouble;
				index += 1;
				continue;
			}
			if (!inSingle && !inDouble) {
				if (ch === "{") {
					braceDepth += 1;
				} else if (ch === "}") {
					braceDepth -= 1;
				} else if (ch === "[") {
					bracketDepth += 1;
				} else if (ch === "]") {
					bracketDepth -= 1;
				} else if (ch === "," && braceDepth === 0 && bracketDepth === 0) {
					break;
				}
			}
			index += 1;
		}

		const value = src.slice(valueStart, index).trim();
		if (key && value) {
			pairs.push({ key, value });
		}
		if (index < src.length && src[index] === ",") {
			index += 1;
		}
	}

	return pairs;
}

function splitTopLevelListItems(listValue: string): string[] {
	const trimmed = listValue.trim();
	if (!(trimmed.startsWith("[") && trimmed.endsWith("]"))) {
		return [];
	}
	const src = trimmed.slice(1, -1).trim();
	if (!src) {
		return [];
	}

	const out: string[] = [];
	let start = 0;
	let index = 0;
	let braceDepth = 0;
	let bracketDepth = 0;
	let inSingle = false;
	let inDouble = false;
	let escaped = false;

	while (index < src.length) {
		const ch = src[index]!;
		if (escaped) {
			escaped = false;
			index += 1;
			continue;
		}
		if (ch === "\\") {
			escaped = true;
			index += 1;
			continue;
		}
		if (ch === "'" && !inDouble) {
			inSingle = !inSingle;
			index += 1;
			continue;
		}
		if (ch === '"' && !inSingle) {
			inDouble = !inDouble;
			index += 1;
			continue;
		}
		if (!inSingle && !inDouble) {
			if (ch === "{") {
				braceDepth += 1;
			} else if (ch === "}") {
				braceDepth -= 1;
			} else if (ch === "[") {
				bracketDepth += 1;
			} else if (ch === "]") {
				bracketDepth -= 1;
			} else if (ch === "," && braceDepth === 0 && bracketDepth === 0) {
				const item = src.slice(start, index).trim();
				if (item) out.push(item);
				start = index + 1;
			}
		}
		index += 1;
	}

	const tail = src.slice(start).trim();
	if (tail) {
		out.push(tail);
	}
	return out;
}

function parseEnchantmentsFromList(value: string): string[] | null {
	const compounds = splitTopLevelListItems(value);
	if (compounds.length === 0) {
		return null;
	}

	const lines: string[] = [];
	for (const compound of compounds) {
		const match = compound.match(ENCHANTMENT_COMPOUND_PATTERN);
		if (!match) {
			return null;
		}
		const enchantmentId = match[1];
		const level = Number.parseInt(match[2] ?? "", 10);
		if (
			!enchantmentId ||
			!Number.isInteger(level) ||
			level < 1 ||
			level > 255
		) {
			return null;
		}
		lines.push(`${enchantmentId} ${level}`);
	}
	return lines;
}

export type ExtractedKnownCustomNbt = {
	enchantmentsFromNbt: string;
	repairCost: string;
	hideFlags: string;
	potionId: string;
	customPotionColor: string;
	customPotionEffects: string;
	attributeModifiers: string;
	remainingCustomNbt: string;
};

export function extractKnownCustomNbt(fragment: string): ExtractedKnownCustomNbt {
	const pairs = splitTopLevelPairs(fragment);
	const remaining: TopLevelPair[] = [];
	let enchantmentsFromNbt = "";
	let repairCost = "";
	let hideFlags = "";
	let potionId = "";
	let customPotionColor = "";
	let customPotionEffects = "";
	let attributeModifiers = "";

	for (const pair of pairs) {
		if (pair.key === "Enchantments") {
			const parsed = parseEnchantmentsFromList(pair.value);
			if (parsed) {
				enchantmentsFromNbt = parsed.join("\n");
				continue;
			}
		}
		if (pair.key === "RepairCost") {
			repairCost = pair.value;
			continue;
		}
		if (pair.key === "HideFlags") {
			hideFlags = pair.value;
			continue;
		}
		if (pair.key === "Potion") {
			potionId = pair.value.replace(/^"/, "").replace(/"$/, "");
			continue;
		}
		if (pair.key === "CustomPotionColor") {
			customPotionColor = pair.value;
			continue;
		}
		if (pair.key === "CustomPotionEffects") {
			customPotionEffects = pair.value;
			continue;
		}
		if (pair.key === "AttributeModifiers") {
			attributeModifiers = pair.value;
			continue;
		}
		remaining.push(pair);
	}

	return {
		enchantmentsFromNbt,
		repairCost,
		hideFlags,
		potionId,
		customPotionColor,
		customPotionEffects,
		attributeModifiers,
		remainingCustomNbt: remaining.map((pair) => `${pair.key}:${pair.value}`).join(","),
	};
}
