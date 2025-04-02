from typing import List
from google.cloud.vision_v1.types.image_annotator import EntityAnnotation


def group_text_by_lines(
    texts: List[EntityAnnotation], y_threshold: int = 10
) -> List[str]:
    """
    Groups text annotations into lines based on their bounding box y-coordinate.

    :param texts: List of EntityAnnotation objects from text_annotations.
    :param y_threshold: Maximum y-difference to consider text on the same line.
    :return: List of lines as strings.
    """
    if not texts:
        return []

    word_annotations = texts[1:]

    words_info = []
    for annotation in word_annotations:
        vertices = annotation.bounding_poly.vertices
        avg_y = sum(vertex.y for vertex in vertices) / len(vertices)
        avg_x = sum(vertex.x for vertex in vertices) / len(vertices)
        words_info.append((annotation.description, avg_y, avg_x))

    # Group words by lines based on y proximity
    lines = []
    for word, y, x in sorted(words_info, key=lambda w: w[1]):  # sort by y
        placed = False
        for line in lines:
            if abs(line[0][1] - y) <= y_threshold:
                line.append((word, y, x))
                placed = True
                break
        if not placed:
            lines.append([(word, y, x)])

    final_lines = []
    for line in lines:
        sorted_line = sorted(line, key=lambda w: w[2])  # sort by x
        line_text = " ".join(word for word, _, _ in sorted_line)
        final_lines.append(line_text)

    return final_lines
