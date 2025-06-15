# 演示

[English Version / 英文版本](README.md)

这个目录包含 Go Playground 项目的演示录制文件和 GIF 图片。

## 录制演示

### 使用 Asciinema

1. 安装 asciinema:

   ```bash
   # macOS
   brew install asciinema

   # Linux
   pip install asciinema
   ```

2. 录制演示:

   ```bash
   # Start recording and run make command directly
   # Note: After the program finishes running, press 'Q' to quit the program and complete the recording
   mkdir -p ./demos/cellular-automaton
   asciinema rec ./demos/cellular-automaton/colorful.cast --title "Cellular Automaton - Colorful" --command "make cellular-automaton-colorful"
   ```

3. 播放演示:

   ```bash
   asciinema play ./demos/cellular-automaton/colorful.cast
   ```

4. 上传到 asciinema.org (可选):

   ```bash
   asciinema upload ./demos/cellular-automaton/colorful.cast
   ```
