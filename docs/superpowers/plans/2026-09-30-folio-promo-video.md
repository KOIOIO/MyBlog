# Folio Promo Video Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Produce a polished approximately 30-second Folio promo video for ordinary visitors, using real site footage, young-neutral narration, synchronized subtitles, and a quiet knowledge-space visual style.

**Architecture:** Keep production assets isolated under `promo-video/`. Capture real browser footage as short clips, create a timed voiceover/subtitle track from one canonical script, then compose the clips, overlays, audio, and end card with FFmpeg. Keep the source timeline and subtitle files editable so a later revision does not require recapturing every shot.

**Tech Stack:** Chrome/Codex browser capture, FFmpeg, WebVTT/ASS subtitles, WAV/MP3 audio, PNG title/end cards.

---

### Task 1: Prepare the production workspace and canonical script

**Files:**
- Create: `promo-video/README.md`
- Create: `promo-video/script.txt`
- Create: `promo-video/shot-list.md`
- Create: `promo-video/raw/`, `promo-video/audio/`, `promo-video/graphics/`, `promo-video/output/`

- [ ] **Step 1: Create the production directories**

Run:
```bash
mkdir -p promo-video/{raw,audio,graphics,output}
```
Expected: all five directories exist.

- [ ] **Step 2: Write the canonical voiceover script**

Create `promo-video/script.txt` with these exact lines and timings:
```text
00:00-00:03  这里，是 Folio。
00:03-00:08  一个可以慢下来，读点东西的地方。
00:08-00:14  技术、生活，还有那些突然冒出来的灵感。
00:14-00:20  用搜索和标签，找到你此刻真正想看的内容。
00:20-00:25  白天专注阅读，夜晚切换成另一种安静。
00:25-00:28  读完，也欢迎留下你的想法。
00:28-00:30  Folio，把思考，留在这里。
```

- [ ] **Step 3: Write the shot list**

Create `promo-video/shot-list.md` mapping each timed line to one real-site capture: title card, homepage/article list, category/tag filtering, search plus article open, article reading plus dark mode, comment/forum/friend-link surfaces, and final brand card with `wwyhahablog.top`.

- [ ] **Step 4: Document capture and export settings**

Create `promo-video/README.md` with the target settings: 1920x1080, 30 fps, H.264 video, AAC audio, 48 kHz mix, 30-second target duration, and mobile-safe margins for all captions.

- [ ] **Step 5: Commit the production scaffold**

Run:
```bash
git add promo-video
git commit -m "chore: scaffold Folio promo video production"
```

### Task 2: Capture real website footage

**Files:**
- Create: `promo-video/raw/01-homepage.mp4`
- Create: `promo-video/raw/02-discover.mp4`
- Create: `promo-video/raw/03-reading.mp4`
- Create: `promo-video/raw/04-community.mp4`

- [ ] **Step 1: Capture the homepage opening**

Open `https://wwyhahablog.top/` at a 16:9 viewport and capture a 5-7 second slow scroll showing Folio branding, article cards, categories, and tags. Do not include the admin account or password in the recording.

- [ ] **Step 2: Capture discovery interactions**

Capture search and category/tag filtering in one 7-9 second clip. Use public content only; keep cursor movement slow and deliberate.

- [ ] **Step 3: Capture reading and dark mode**

Open a public article and capture 6-8 seconds of reading, then transition to dark mode for the final 2 seconds.

- [ ] **Step 4: Capture community surfaces**

Capture 3-4 seconds showing comments, forum, friend links, or contact entry points without exposing private user information.

- [ ] **Step 5: Verify each clip**

Run:
```bash
ffprobe -v error -show_entries stream=width,height,r_frame_rate,duration -of default=noprint_wrappers=1 promo-video/raw/*.mp4
```
Expected: every clip is readable, 16:9, approximately 30 fps, and contains no credential or private-data exposure.

### Task 3: Create voiceover, subtitles, and graphics

**Files:**
- Create: `promo-video/audio/voiceover.wav`
- Create: `promo-video/audio/music.wav`
- Create: `promo-video/subtitles.ass`
- Create: `promo-video/graphics/title.png`
- Create: `promo-video/graphics/end-card.png`

- [ ] **Step 1: Record or synthesize the young-neutral voiceover**

Produce one clean mono voice track matching `promo-video/script.txt`, with natural pauses and total duration between 28 and 30 seconds. Normalize peaks below -1 dBFS.

- [ ] **Step 2: Prepare background music**

Use a lyric-free lo-fi piano or acoustic-guitar bed. Trim or loop it to 30 seconds and duck it under speech to approximately -24 LUFS integrated while speech remains intelligible.

- [ ] **Step 3: Author timed subtitles**

Create `promo-video/subtitles.ass` with one event per script line, 8-14 Chinese characters per line where possible, bottom safe-zone placement, and warm-white text with a dark translucent outline.

- [ ] **Step 4: Create title and end-card graphics**

Create 1920x1080 PNGs using warm white, ink green, and restrained blue-gray. The title reads `Folio` and `把思考，留在这里。`; the end card reads `Folio`, `把思考，留在这里。`, and `wwyhahablog.top`.

- [ ] **Step 5: Inspect audio and subtitle timing**

Run:
```bash
ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1 promo-video/audio/voiceover.wav
```
Expected: duration is within 28-30 seconds; subtitle events cover the voiceover with no overlap or gap that changes meaning.

### Task 4: Assemble and export the promo video

**Files:**
- Create: `promo-video/output/folio-promo-30s.mp4`
- Create: `promo-video/output/folio-promo-30s-muted.mp4`

- [ ] **Step 1: Build the visual timeline**

Concatenate or trim the four clips to the shot-list durations, add title and end cards, apply gentle crossfades, and preserve enough headroom for subtitles.

- [ ] **Step 2: Mix audio and burn subtitles**

Mix normalized voiceover with ducked music, burn `subtitles.ass`, and export H.264/AAC at 1920x1080, 30 fps, `yuv420p`, and `+faststart`.

- [ ] **Step 3: Export the muted accessibility preview**

Create the same video without an audio track so subtitle readability can be checked independently.

- [ ] **Step 4: Verify the finished files**

Run:
```bash
ffprobe -v error -show_entries format=duration:stream=codec_name,width,height,r_frame_rate,pix_fmt -of default=noprint_wrappers=1 promo-video/output/folio-promo-30s.mp4
```
Expected: approximately 30 seconds, H.264 video, AAC audio, 1920x1080, 30 fps, `yuv420p`.

### Task 5: Visual and content QA

- [ ] **Step 1: Render contact frames**

Run:
```bash
mkdir -p promo-video/output/contact-sheet
ffmpeg -y -i promo-video/output/folio-promo-30s.mp4 -vf "fps=1,scale=480:-1,tile=5x2" promo-video/output/contact-sheet/contact-sheet.jpg
```
Expected: a 10-frame contact sheet showing all major beats.

- [ ] **Step 2: Check the review criteria**

Confirm: Folio and URL are legible; website features shown are real; no credentials/private data appear; captions never cover important UI; narration matches captions; total duration is about 30 seconds; muted playback remains understandable.

- [ ] **Step 3: Commit the final production artifacts**

Run:
```bash
git add promo-video
git commit -m "feat: produce Folio 30-second promo video"
```
