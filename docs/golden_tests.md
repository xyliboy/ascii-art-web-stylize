# Golden Tests

## Notes
- \n in input represents a real newline character (Enter key in the web form)
- Each character is rendered across 8 rows; the 8th row is often blank (spaces only)
- Trailing spaces on each line are normalized (trimmed) before comparison
- Multi-line input produces segments separated by two blank lines (8th row + separator)

## Test Case 1
Banner: standard
Input:
```
{123}
<Hello> (World)!
```
Expected Output:
```
   __                     __
  / /  _   ____    _____  \ \
 | |  / | |___ \  |___ /   | |
/ /   | |   __) |   |_ \    \ \
\ \   | |  / __/   ___) |   / /
 | |  |_| |_____| |____/   | |
  \_\                     /_/


   __  _    _          _   _          __            __ __          __                 _       _  __    _
  / / | |  | |        | | | |         \ \          / / \ \        / /                | |     | | \ \  | |
 / /  | |__| |   ___  | | | |   ___    \ \        | |   \ \  /\  / /    ___    _ __  | |   __| |  | | | |
< <   |  __  |  / _ \ | | | |  / _ \    > >       | |    \ \/  \/ /    / _ \  | '__| | |  / _` |  | | | |
 \ \  | |  | | |  __/ | | | | | (_) |  / /        | |     \  /\  /    | (_) | | |    | | | (_| |  | | |_|
  \_\ |_|  |_|  \___| |_| |_|  \___/  /_/         | |      \/  \/      \___/  |_|    |_|  \__,_|  | | (_)
                                                   \_\                                           /_/
```
Two blank lines between segments: 8th row of segment 1 + separator line.

## Test Case 2
Banner: standard
Input:
```
123??
```
Expected Output:
```
                     ___    ___
 _   ____    _____  |__ \  |__ \
/ | |___ \  |___ /     ) |    ) |
| |   __) |   |_ \    / /    / /
| |  / __/   ___) |  |_|    |_|
|_| |_____| |____/   (_)    (_)
```

## Test Case 3
Banner: shadow
Input:
```
$% "=
```
Expected Output:
```
                        _|  _|
  _|   _|_|    _|       _|  _|
_|_|_| _|_|  _|                _|_|_|_|_|
_|_|       _|
  _|_|   _|  _|_|              _|_|_|_|_|
_|_|_| _|    _|_|
  _|
```

## Test Case 4
Banner: thinkertoy
Input:
```
123 T/fs#R
```
Expected Output:
```

  0    --  o-o        o-O-o     o  o-o      | |  o--o
 /|   o  o    |         |      /   |       -O-O- |   |
o |     /   oo          |     o   -O-  o-o  | |  O-Oo
  |    /      |         |    /     |    \  -O-O- |  \
o-o-o o--o o-o          o   o      o   o-o  | |  o   o
```
Starts with a blank first row — the thinkertoy banner's row 1 for these characters is empty.
