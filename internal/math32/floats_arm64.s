//go:build !noasm && arm64

TEXT ·vdotNEON(SB), $0-32
	MOVD a+0(FP), R0          // Move the value of 'a' into register R0
	MOVD b+8(FP), R1          // Move the value of 'b' into register R1
	MOVD n+16(FP), R2         // Move the value of 'n' into register R2
	MOVD ret+24(FP), R3       // Move the address of the return value into register R3
	WORD $0xa9bf7bfd          // Save the frame pointer and link register to the stack
	WORD $0x91000c48          // Add the value of register x2 to register x8 and store the result in x8
	WORD $0xf100005f          // Compare the value in register x2 with 0 and set flags
	WORD $0x9a82b108          // Conditional select: If the previous comparison result is less than, set x8 to x2, else keep x8 unchanged
	WORD $0x9342fd0a          // Arithmetic shift right: Shift the value in x8 right by the number of bits specified in x10 and store the result in x10
	WORD $0x927ef508          // Bitwise AND: Perform a bitwise AND operation between the values in x8 and x8, store the result in x8
	WORD $0x7100055f          // Compare the value in w10 with 0 and set flags
	WORD $0xcb080048          // Subtract the value in x8 from the value in x2 and store the result in x8
	WORD $0x910003fd          // Move the value in the stack pointer to register x29
	WORD $0x540002ab          // Branch to label .LBB4_5 if the previous comparison result is less than
	WORD $0x3cc10400          // Load a quadword from the memory address stored in register x0 into vector register q0
	WORD $0x3cc10421          // Load a quadword from the memory address stored in register x1 into vector register q1
	WORD $0x71000549          // Subtract the value in w10 from w9 and set flags
	WORD $0x6e21dc00          // Multiply the vectors in v0 and v1 element-wise and store the result in v0
	WORD $0x54000200          // Branch to label .LBB4_6 if the previous comparison result is equal to
	WORD $0xb27d7beb          // Move the value in register x11 to register x11
	WORD $0x8b0a096a          // Add the value in x11 shifted left by the value in x10 to the value in x10 and store the result in x10
	WORD $0x927e7d4a          // Bitwise AND: Perform a bitwise AND operation between the values in x10 and x10, store the result in x10
	WORD $0x9100114b          // Add the value in x11 to the value in x10 and store the result in x11
	WORD $0x8b0b080a          // Add the value in x0 to the value in x11 shifted left by the value in x10 and store the result in x10
	WORD $0xaa0103ec          // Move the value in register x1 to register x12

LBB4_3:
	WORD $0x3cc10401          // Load a quadword from the memory address stored in register x0 into vector register q1
	WORD $0x3cc10582          // Load a quadword from the memory address stored in register x12 into vector register q2
	WORD $0x71000529          // Subtract the value in w10 from w9 and set flags
	WORD $0x6e22dc21          // Multiply the vectors in v1 and v2 element-wise and store the result in v1
	WORD $0x4e21d400          // Add the vectors in v0 and v1 and store the result in v0
	WORD $0x54ffff61          // Branch to label .LBB4_3 if the previous comparison result is not equal to
	WORD $0x8b0b0821          // Add the value in x11 to the value in x1 and store the result in x1
	WORD $0xaa0a03e0          // Move the value in register x10 to register x0
	WORD $0x14000001          // Unconditional branch to label .LBB4_6

LBB4_5:
LBB4_6:
	WORD $0x1e2703e1          // Move zero to register s1
	WORD $0x5e0c0402          // Move the value in register v0.s[1] to register s2
	WORD $0x5e140403          // Move the value in register v0.s[2] to register s3
	WORD $0x5e1c0404          // Move the value in register v0.s[3] to register s4
	WORD $0x1e212800          // Add the value in s0 to the value in s1 and store the result in s0
	WORD $0x1e202840          // Add the value in s2 to the value in s0 and store the result in s0
	WORD $0x1e202860          // Add the value in s3 to the value in s0 and store the result in s0
	WORD $0x1e202880          // Add the value in s4 to the value in s0 and store the result in s0
	WORD $0x7100011f          // Compare the value in w8 with 0 and set flags
	WORD $0xbd000060          // Store the value in register s0 to the memory address stored in register x3
	WORD $0x5400012d          // Branch to label .LBB4_9 if the previous comparison result is less than or equal to
	WORD $0x92407d08          // Bitwise AND: Perform a bitwise AND operation between the values in x8 and x8, store the result in x8

LBB4_8:
	WORD $0xbc404401          // Load a single precision floating-point value from the memory address stored in register x0 into register s1
	WORD $0xbc404422          // Load a single precision floating-point value from the memory address stored in register x1 into register s2
	WORD $0xf1000508          // Subtract the value in x8 from 0 and set flags
	WORD $0x1e220821          // Multiply the values in s1 and s2 and store the result in s1
	WORD $0x1e212800          // Add the value in s0 to the value in s1 and store the result in s0
	WORD $0xbd000060          // Store the value in register s0 to the memory address stored in register x3
	WORD $0x54ffff41          // Branch to label .LBB4_8 if the previous comparison result is not equal to

LBB4_9:
	WORD $0xa8c17bfd          // Load the frame pointer and link register from the stack
	WORD $0xd65f03c0          // Return
