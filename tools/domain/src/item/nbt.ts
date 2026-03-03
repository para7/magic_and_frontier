import type { SaveItemInput } from "./types.js";

export type BuildItemNbtInput = SaveItemInput;

export type NbtBuildResult = {
	nbt: string;
	enchantmentError?: string;
};

function toNbtText(value: string): string {
	const json = JSON.stringify({ text: value }).replace(/'/g, "\\'");
	return `'${json}'`;
}

function normalizeCustomNbtFragment(value: string): string {
	const trimmed = value.trim();
	if (!trimmed) {
		return "";
	}
	if (trimmed.startsWith("{") && trimmed.endsWith("}")) {
		return trimmed.slice(1, -1).trim();
	}
	return trimmed;
}

function normalizeListValue(value: string): string {
	const trimmed = value.trim();
	if (!trimmed) {
		return "";
	}
	if (trimmed.startsWith("[") && trimmed.endsWith("]")) {
		return trimmed;
	}
	return `[${trimmed}]`;
}

function normalizePotionId(value: string): string {
	return value.trim().replace(/^"/, "").replace(/"$/, "");
}

export function buildItemNbt(input: BuildItemNbtInput): NbtBuildResult {
	const tagParts: string[] = [];

	const displayParts: string[] = [];
	if (input.customName.trim().length > 0) {
		displayParts.push(`Name:${toNbtText(input.customName.trim())}`);
	}

	const loreLines = input.lore
		.split(/\r?\n/)
		.map((line) => line.trim())
		.filter(Boolean);
	if (loreLines.length > 0) {
		const lore = loreLines.map((line) => toNbtText(line)).join(",");
		displayParts.push(`Lore:[${lore}]`);
	}

	if (displayParts.length > 0) {
		tagParts.push(`display:{${displayParts.join(",")}}`);
	}

	const enchantmentEntries = input.enchantments
		.split(/\r?\n/)
		.map((line) => line.trim())
		.filter(Boolean);
	if (enchantmentEntries.length > 0) {
		const enchantments: string[] = [];
		for (const entry of enchantmentEntries) {
			const [rawId, rawLevel] = entry.split(/\s+/, 2);
			const level = Number.parseInt(rawLevel ?? "", 10);
			if (!rawId || !Number.isInteger(level) || level < 1 || level > 255) {
				return {
					nbt: "",
					enchantmentError: `Invalid enchantment line: "${entry}". Use "minecraft:sharpness 5" format.`,
				};
			}
			enchantments.push(`{id:"${rawId}",lvl:${level}s}`);
		}
		tagParts.push(`Enchantments:[${enchantments.join(",")}]`);
	}

	if (input.unbreakable) {
		tagParts.push("Unbreakable:1b");
	}

	if (input.customModelData.trim().length > 0) {
		tagParts.push(`CustomModelData:${input.customModelData.trim()}`);
	}

	if (input.repairCost.trim().length > 0) {
		tagParts.push(`RepairCost:${input.repairCost.trim()}`);
	}

	if (input.hideFlags.trim().length > 0) {
		tagParts.push(`HideFlags:${input.hideFlags.trim()}`);
	}

	const potionId = normalizePotionId(input.potionId);
	if (potionId.length > 0) {
		tagParts.push(`Potion:"${potionId}"`);
	}

	if (input.customPotionColor.trim().length > 0) {
		tagParts.push(`CustomPotionColor:${input.customPotionColor.trim()}`);
	}

	const customPotionEffects = normalizeListValue(input.customPotionEffects);
	if (customPotionEffects.length > 0) {
		tagParts.push(`CustomPotionEffects:${customPotionEffects}`);
	}

	const attributeModifiers = normalizeListValue(input.attributeModifiers);
	if (attributeModifiers.length > 0) {
		tagParts.push(`AttributeModifiers:${attributeModifiers}`);
	}

	const customFragment = normalizeCustomNbtFragment(input.customNbt);
	if (customFragment.length > 0) {
		tagParts.push(customFragment);
	}

	const itemNbtParts = [`id:"${input.itemId.trim()}"`, `Count:${input.count}b`];
	if (tagParts.length > 0) {
		itemNbtParts.push(`tag:{${tagParts.join(",")}}`);
	}

	return {
		nbt: `{${itemNbtParts.join(",")}}`,
	};
}
