#!/usr/bin/env python3
import json


def extract_entries(item):
    entries = []

    names = item.get("names")
    emoji = item.get("surrogates")
    if names and emoji:
        entries.append({"name": names[0], "emoji": emoji})

    for child in item.get("diversityChildren", []):
        child_names = child.get("names")
        child_emoji = child.get("surrogates")
        if child_names and child_emoji:
            entries.append({"name": child_names[0], "emoji": child_emoji})

    return entries


def simplify_all_categories(input_path, output_path):
    with open(input_path, "r", encoding="utf-8") as f:
        data = json.load(f)

    result = []

    for category, items in data.items():
        if isinstance(items, list):
            for item in items:
                result.extend(extract_entries(item))

    with open(output_path, "w", encoding="utf-8") as f:
        json.dump(result, f, ensure_ascii=False, indent=2)


if __name__ == "__main__":
    simplify_all_categories("discord-emojis.pretty.json", "discord-emojis.flat.json")
