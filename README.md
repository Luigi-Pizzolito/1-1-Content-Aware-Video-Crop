# 1:1 Content-Aware Video Crop
 
This tool is designed to crop videos into a square format optimized for hologram projection. It uses content-aware cropping to ensure that important content remains visible after cropping.

<details open>
<summary>Sample Output</summary>

![vidout3](https://github.com/Luigi-Pizzolito/1-1-Content-Aware-Video-Crop/assets/27804554/e3e42d7a-b630-40aa-85e8-04b732b6dbe0)

</details>

## Usage

### Command-Line Flags

The tool accepts the following command-line flags:

- `-i`: A string flag indicating the input video file(s) to be cropped. Multiple input videos can be provided by providing a folder instead. No default value; this flag is required.
- `-o`: A string flag indicating the output directory for the cropped videos. Default is current working directory.
- `-s`: An integer flag indicating the size of the square output video. Default is `256`.
- `--play`: A boolean flag indicating whether to enable player-only mode, does not save videos, just displays a preview of the result. Default is `false`.
- `--ui`: A boolean flag indicating whether to draw a user interface during processing. Default is `false`.
- `--rt`: A boolean flag indicating whether to process the video in real-time. Default is `false`.

### Example Usage
<details open>
<summary>Algorithm Display Mode</summary>

```bash
go run . -i input.mp4 -o output_dir -s 256 --ui --rt
```

![show](https://github.com/Luigi-Pizzolito/1-1-Content-Aware-Video-Crop/assets/27804554/fb044c77-3723-46e8-8f88-102468bfb0f5)
</details>


<details>
<summary>Player-Only Mode</summary>

```bash
go run . -i input.mp4 -s 256 --play
```

![play](https://github.com/Luigi-Pizzolito/1-1-Content-Aware-Video-Crop/assets/27804554/01ca8b2d-2d0f-4bde-99d9-b72e7d8e6421)
</details>


<details>
<summary>Headless Mode (default)</summary>

```bash
go run . -i input.mp4 -o output_dir -s 256
```

![cmd](https://github.com/Luigi-Pizzolito/1-1-Content-Aware-Video-Crop/assets/27804554/a5ee6b02-2740-4388-88af-154af6f8a682)
</details>


<details>
<summary>Providing folder of videos as input</summary>

```bash
go run . -i input_dir -o output_dir
```
</details>