# Little Man Computer (LMC) Interpreter in Go

An extended **Little Man Computer (LMC)** interpreter and assembler built from scratch in Go. This tool reads standard LMC assembly code, compiles it into a 3-digit decimal instruction array (simulated RAM), and executes it using a simulated CPU cycle.

In addition to standard LMC commands, this implementation includes the **`CHAR` (OUTCHAR)** extension, allowing programs to output text characters and print words or full sentences.

---

## Features
- **Two-Pass Assembler:** Automatically resolves symbol labels and variables into 00-99 memory addresses.
- **Classic 3-Digit Architecture:** Accurately simulates the Accumulator, Program Counter (PC), and 100 memory cells.
- **`CHAR` Extension:** Built-in `OUTCHAR` function (Opcode `903`) converts decimal accumulator values to ASCII/UTF-8 text characters on the fly.
- **Error Handling:** Helpful assembly-time and runtime error messaging for bad mnemonics, missing labels, and infinite loops.

---

## Instructions

| Mnemonic | Numeric Opcode | Description |
| :--- | :--- | :--- |
| **ADD** | `1xx` | Add the contents of address `xx` to the Accumulator. |
| **SUB** | `2xx` | Subtract the contents of address `xx` from the Accumulator. |
| **STA** / **STO** | `3xx` | Store the contents of the Accumulator at address `xx`. |
| **LDA** | `5xx` | Load the contents of address `xx` into the Accumulator. |
| **BRA** | `6xx` | Branch unconditionally to address `xx` (sets PC to `xx`). |
| **BRZ** | `7xx` | Branch to address `xx` if the Accumulator value is exactly `0`. |
| **BRP** | `8xx` | Branch to address `xx` if the Accumulator is `0` or positive. |
| **INP** / **IN** | `901` | Prompt user for a numeric integer input and save to the Accumulator. |
| **OUT** | `902` | Print the current numeric value of the Accumulator to the screen. |
| **HLT** / **COB** | `000` | Halt execution / Coffee Break. |
| **DAT** | *None* | Reserves a memory address for data storage. Can be optionally pre-initialized. |
| **CHAR** | `903` | **EXTRA INSTRUCTION - NOT IN ACTUAL LMC ASSEMBLY** Treats the value in the Accumulator as an ASCII code and prints its character (e.g., `65` -> `'A'`). Does not append a newline. |

---

