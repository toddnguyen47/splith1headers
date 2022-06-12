"""Split h1 headers"""

import os
from typing import List, Dict

_FILE_PATH = "Z:/DocumentsAndStuff/Desktop/1.xhtml"
_OUTFILE = "Body"

_FILE_PATH_DIRNAME = os.path.dirname(_FILE_PATH)
_, _EXTENSION = os.path.splitext(_FILE_PATH)

_IMAGE_XML = r'<image height="952" width="646" xlink:href='
_IMAGE_SRC = r"<img src="
_CITE_NOTE = r'href="#cite_note'


def _main():
    """Main Function to execute"""

    with open(_FILE_PATH, "r", encoding="utf-8") as file:
        data = file.read()

    notes_section_exists = False
    # Find h1 lines
    h1_lines = []
    lines = data.splitlines()
    for index, line in enumerate(lines):
        if r"<h1" in line:
            h1_lines.append(index)
            if "notes" in line.lower():
                notes_section_exists = True

    # Split into files
    contents: List[List[str]] = []
    for i, _ in enumerate(h1_lines):
        if i == len(h1_lines) - 1:
            current_lines = lines[h1_lines[i] :]
        else:
            current_lines = lines[h1_lines[i] : h1_lines[i + 1]]
        contents.append(current_lines)

    cite_ref_reverse_map: Dict[str, int] = {}
    # Per file
    for index, content in enumerate(contents):
        for index2, line in enumerate(content):
            if (
                r'<div class="svg_outer svg_inner"' in line
                or r'<svg xmlns="http://www.w3.org/2000/svg" height="100%"' in line
                or r"</svg>" in line
                or r"</div>" in line
            ):
                continue
            if _IMAGE_XML in line:
                line = line.replace(_IMAGE_XML, _IMAGE_SRC)
                line = line.replace(r"</image>", "</img>")
            content[index2] = line

            # Replace cite_ref
            if _CITE_NOTE in line:
                line = line.replace(_CITE_NOTE, r'href="Notes.xhtml#cite_note')

        contents[index] = content

    # Output
    # for i, _ in enumerate(h1_lines):
    #     # If last file
    #     outfile_name = os.path.join(_FILE_PATH_DIRNAME, _OUTFILE + str(i) + _EXTENSION)
    #     if i == len(h1_lines) - 1:
    #         current_lines = lines[h1_lines[i] :]
    #     else:
    #         current_lines = lines[h1_lines[i] : h1_lines[i + 1]]
    #     print(f"Writing to {outfile_name}")
    #     with open(outfile_name, "w", encoding="utf-8") as outfile:
    #         outfile.write("\n".join(current_lines))


if __name__ == "__main__":
    _main()
