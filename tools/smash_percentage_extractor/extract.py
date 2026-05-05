import cv2
import sys
import re
from pathlib import Path


def extract_percentages(video_path, output_dir="output", interval=5):
    video_path = Path(video_path)
    output_dir = Path(output_dir)
    output_dir.mkdir(exist_ok=True)

    cap = cv2.VideoCapture(str(video_path))
    fps = cap.get(cv2.CAP_PROP_FPS)
    frame_interval = int(fps * interval)

    width = int(cap.get(cv2.CAP_PROP_FRAME_WIDTH))
    height = int(cap.get(cv2.CAP_PROP_FRAME_HEIGHT))

    # x1, y1, x2, y2
    p1_region = (
        int(width * 0.26),
        int(height * 0.842),
        int(width * 0.3515),
        int(height * 0.932),
    )
    p2_region = (
        int(width * 0.64),
        int(height * 0.842),
        int(width * 0.737),
        int(height * 0.932),
    )

    frame_count = 0
    base_name = re.sub(r"[^A-Za-z0-9]", "", video_path.stem)[:20]

    while True:
        ret, frame = cap.read()
        if not ret:
            break

        if frame_count % frame_interval == 0:
            time_sec = int(frame_count / fps)

            # Extract P1 percentage
            x1, y1, x2, y2 = p1_region
            p1_crop = frame[y1:y2, x1:x2]
            p1_filename = f"{base_name}_p1_{time_sec:04d}s.png"
            cv2.imwrite(str(output_dir / p1_filename), p1_crop)

            # Extract P2 percentage
            x1, y1, x2, y2 = p2_region
            p2_crop = frame[y1:y2, x1:x2]
            p2_filename = f"{base_name}_p2_{time_sec:04d}s.png"
            cv2.imwrite(str(output_dir / p2_filename), p2_crop)

            print(f"Extracted frame at {time_sec}s")

        frame_count += 1
        if frame_count // frame_interval > 100:
            break

    cap.release()
    print(f"Done. Extracted {frame_count // frame_interval} frames per player.")


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print(
            "Usage: python extract_percentages.py <video_file.mp4> [output_dir] [interval_seconds]"
        )
        sys.exit(1)

    video = sys.argv[1]
    output = sys.argv[2] if len(sys.argv) > 2 else "output"
    interval = int(sys.argv[3]) if len(sys.argv) > 3 else 5

    extract_percentages(video, output, interval)
