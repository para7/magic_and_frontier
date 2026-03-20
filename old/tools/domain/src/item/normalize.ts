import { buildItemNbt } from "./nbt.js";
import { extractKnownCustomNbt } from "./custom-nbt.js";
import { defaultItemState, type ItemEntry, type ItemState } from "./types.js";

function fallbackNbtFor(record: Record<string, unknown>): string {
	const built = buildItemNbt({
		id: "fallback",
		itemId:
			typeof record.itemId === "string" ? record.itemId : "minecraft:stone",
		count:
			typeof record.count === "number" && Number.isFinite(record.count)
				? Math.max(1, Math.floor(record.count))
				: 1,
		customName: typeof record.customName === "string" ? record.customName : "",
		lore: typeof record.lore === "string" ? record.lore : "",
		enchantments:
			typeof record.enchantments === "string" ? record.enchantments : "",
		unbreakable:
			typeof record.unbreakable === "boolean" ? record.unbreakable : false,
		customModelData:
			typeof record.customModelData === "string" ? record.customModelData : "",
		repairCost:
			typeof record.repairCost === "string" ? record.repairCost : "",
		hideFlags: typeof record.hideFlags === "string" ? record.hideFlags : "",
		potionId: typeof record.potionId === "string" ? record.potionId : "",
		customPotionColor:
			typeof record.customPotionColor === "string"
				? record.customPotionColor
				: "",
		customPotionEffects:
			typeof record.customPotionEffects === "string"
				? record.customPotionEffects
				: "",
		attributeModifiers:
			typeof record.attributeModifiers === "string"
				? record.attributeModifiers
				: "",
		customNbt: typeof record.customNbt === "string" ? record.customNbt : "",
	});
	if (!built.enchantmentError) {
		return built.nbt;
	}
	const safeItemId =
		typeof record.itemId === "string" && record.itemId.length > 0
			? record.itemId
			: "minecraft:stone";
	const safeCount =
		typeof record.count === "number" && Number.isFinite(record.count)
			? Math.max(1, Math.floor(record.count))
			: 1;
	return `{id:"${safeItemId}",Count:${safeCount}b}`;
}

function normalizeItemEntry(value: unknown): ItemEntry | null {
	if (!value || typeof value !== "object") {
		return null;
	}

	const record = value as Record<string, unknown>;
	const extracted = extractKnownCustomNbt(
		typeof record.customNbt === "string" ? record.customNbt : "",
	);
	const fallbackNbt = fallbackNbtFor(record);

	return {
		id:
			typeof record.id === "string" && record.id.length > 0
				? record.id
				: crypto.randomUUID(),
		itemId:
			typeof record.itemId === "string" && record.itemId.length > 0
				? record.itemId
				: "minecraft:stone",
		count:
			typeof record.count === "number" && Number.isFinite(record.count)
				? Math.max(1, Math.floor(record.count))
				: 1,
		customName: typeof record.customName === "string" ? record.customName : "",
		lore: typeof record.lore === "string" ? record.lore : "",
		enchantments:
			typeof record.enchantments === "string" && record.enchantments.length > 0
				? record.enchantments
				: extracted.enchantmentsFromNbt,
		unbreakable:
			typeof record.unbreakable === "boolean" ? record.unbreakable : false,
		customModelData:
			typeof record.customModelData === "string" ? record.customModelData : "",
		repairCost:
			typeof record.repairCost === "string"
				? record.repairCost
				: extracted.repairCost,
		hideFlags:
			typeof record.hideFlags === "string" ? record.hideFlags : extracted.hideFlags,
		potionId:
			typeof record.potionId === "string" ? record.potionId : extracted.potionId,
		customPotionColor:
			typeof record.customPotionColor === "string"
				? record.customPotionColor
				: extracted.customPotionColor,
		customPotionEffects:
			typeof record.customPotionEffects === "string"
				? record.customPotionEffects
				: extracted.customPotionEffects,
		attributeModifiers:
			typeof record.attributeModifiers === "string"
				? record.attributeModifiers
				: extracted.attributeModifiers,
		customNbt: extracted.remainingCustomNbt,
		nbt:
			typeof record.nbt === "string" && record.nbt.length > 0
				? record.nbt
				: fallbackNbt,
		updatedAt: typeof record.updatedAt === "string" ? record.updatedAt : "",
	};
}

export function normalizeItemState(value: unknown): ItemState {
	if (!value || typeof value !== "object") {
		return defaultItemState;
	}

	const record = value as Record<string, unknown>;
	if (!Array.isArray(record.items)) {
		return defaultItemState;
	}

	const items = record.items
		.map(normalizeItemEntry)
		.filter((item): item is ItemEntry => item !== null);
	return { items };
}
