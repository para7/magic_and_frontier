# Frontend (SSR) Migration Runbook

## Objective
- Replace Angular frontend with Go SSR (`templ + htmx`) while keeping server API contracts available.

## Input References
- `../tools/frontend/src/app/**`
- Existing `tools2` templ usage:
1. `app/views/layout.templ`
2. `app/views/form/*.templ`

## Output Location
- Recommended package roots:
1. `app/internal/web`
2. `app/views/layout*`
3. `app/views/<feature>/*`

## Task List
1. Define navigation and page shells
- Sidebar/nav entries for six features.
- Shared layout template with message/toast area.
- Active-route highlighting.

2. Implement feature pages in migration order
1. `items`
2. `grimoire`
3. `skills`
4. `enemy-skills`
5. `treasures`
6. `enemies`

3. Implement per-feature SSR interactions
- List current entries.
- Create/update form submit.
- Delete action.
- htmx partial rerender targets for list + form errors.

4. Implement shared UI fragments
- Field error renderer.
- Reference picker/select component.
- JSON/textarea helper for complex fields.

5. Implement save flow
- Global "save/export" button in layout.
- POST to `/api/save`, show success/failure toast.
- Non-blocking UI update after save completion.

6. Route wiring
- Web routes under non-API paths (for example `/items`, `/grimoire`, ...).
- Keep `/api/*` routes intact for compatibility.

7. Remove/replace old demo form
- Remove `form` feature templates and handlers once new pages are reachable.
- Update root redirect to first migrated feature page.

8. Add UI smoke tests (happy path)
- Page load for each feature.
- Create then list visible for each feature.
- Delete removes entry from list.
- Save button triggers success result display.

## UX Constraints
- Keep implementation utilitarian and consistent with current `tools2` style.
- Favor clear validation feedback over visual polish for this phase.
- Avoid client-heavy JS; htmx snippets only.

## Completion Checklist
- [ ] All six feature screens functional.
- [ ] Save/export button functional from UI.
- [ ] No Angular runtime dependency remains in `tools2`.
- [ ] `templ generate` and `go test ./...` pass.
