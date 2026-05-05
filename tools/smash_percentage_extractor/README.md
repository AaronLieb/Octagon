# Smash Percentage Extractor

Extracts player percentage regions from Super Smash Bros Ultimate videos for training data collection.

## Setup

```bash
pip install -r requirements.txt
```

## Usage

```bash
python extract.py <video_file.mp4> [output_dir] [interval_seconds]
```

**Examples:**

```bash
python extract.py match1.mp4
python extract.py match2.mp4 custom_output 10
```

## Output

Files are named: `{video_name}_p1_{time}s.png` and `{video_name}_p2_{time}s.png`

## Note

You may need to adjust the percentage region coordinates in the script based on your video footage.
