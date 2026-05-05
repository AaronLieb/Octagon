import cv2
import sys
from pathlib import Path


def preprocess_image(input_path, output_path, target_size=(64, 32)):
    img = cv2.imread(str(input_path))
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    resized = cv2.resize(gray, target_size)
    cv2.imwrite(str(output_path), resized)


if __name__ == "__main__":
    input_dir = Path(sys.argv[1]) if len(sys.argv) > 1 else Path("output")
    output_dir = Path(sys.argv[2]) if len(sys.argv) > 2 else Path("preprocessed")
    target_size = (
        (int(sys.argv[3]), int(sys.argv[4])) if len(sys.argv) > 4 else (64, 32)
    )

    output_dir.mkdir(exist_ok=True)

    for img_path in input_dir.glob("*.png"):
        preprocess_image(img_path, output_dir / img_path.name, target_size)
        print(f"Processed: {img_path.name}")
