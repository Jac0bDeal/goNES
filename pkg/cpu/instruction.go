package cpu

type instruction struct {
	name        string
	operate     func() uint8
	addressMode func()
	cycles      uint8
}

// placing this here to prevent editor from lagging due to file size.
func buildLookupTable(cpu Mos6502) [256]instruction {
	return [256]instruction{
		{"BRK", cpu.brk, cpu.imm, 7}, {"ORA", cpu.ora, cpu.izx, 6}, {"???", cpu.xxx, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 8}, {"???", cpu.nop, cpu.imp, 3}, {"ORA", cpu.ora, cpu.zp0, 3}, {"ASL", cpu.asl, cpu.zp0, 5}, {"???", cpu.xxx, cpu.imp, 5}, {"PHP", cpu.php, cpu.imp, 3}, {"ORA", cpu.ora, cpu.imm, 2}, {"ASL", cpu.asl, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 2}, {"???", cpu.nop, cpu.imp, 4}, {"ORA", cpu.ora, cpu.abs, 4}, {"ASL", cpu.asl, cpu.abs, 6}, {"???", cpu.xxx, cpu.imp, 6},
		{"BPL", cpu.bpl, cpu.rel, 2}, {"ORA", cpu.ora, cpu.izy, 5}, {"???", cpu.xxx, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 8}, {"???", cpu.nop, cpu.imp, 4}, {"ORA", cpu.ora, cpu.zpx, 4}, {"ASL", cpu.asl, cpu.zpx, 6}, {"???", cpu.xxx, cpu.imp, 6}, {"CLC", cpu.clc, cpu.imp, 2}, {"ORA", cpu.ora, cpu.aby, 4}, {"???", cpu.nop, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 7}, {"???", cpu.nop, cpu.imp, 4}, {"ORA", cpu.ora, cpu.abx, 4}, {"ASL", cpu.asl, cpu.abx, 7}, {"???", cpu.xxx, cpu.imp, 7},
		{"JSR", cpu.jsr, cpu.abs, 6}, {"AND", cpu.and, cpu.izx, 6}, {"???", cpu.xxx, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 8}, {"BIT", cpu.bit, cpu.zp0, 3}, {"AND", cpu.and, cpu.zp0, 3}, {"ROL", cpu.rol, cpu.zp0, 5}, {"???", cpu.xxx, cpu.imp, 5}, {"PLP", cpu.plp, cpu.imp, 4}, {"AND", cpu.and, cpu.imm, 2}, {"ROL", cpu.rol, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 2}, {"BIT", cpu.bit, cpu.abs, 4}, {"AND", cpu.and, cpu.abs, 4}, {"ROL", cpu.rol, cpu.abs, 6}, {"???", cpu.xxx, cpu.imp, 6},
		{"BMI", cpu.bmi, cpu.rel, 2}, {"AND", cpu.and, cpu.izy, 5}, {"???", cpu.xxx, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 8}, {"???", cpu.nop, cpu.imp, 4}, {"AND", cpu.and, cpu.zpx, 4}, {"ROL", cpu.rol, cpu.zpx, 6}, {"???", cpu.xxx, cpu.imp, 6}, {"SEC", cpu.sec, cpu.imp, 2}, {"AND", cpu.and, cpu.aby, 4}, {"???", cpu.nop, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 7}, {"???", cpu.nop, cpu.imp, 4}, {"AND", cpu.and, cpu.abx, 4}, {"ROL", cpu.rol, cpu.abx, 7}, {"???", cpu.xxx, cpu.imp, 7},
		{"RTI", cpu.rti, cpu.imp, 6}, {"EOR", cpu.eor, cpu.izx, 6}, {"???", cpu.xxx, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 8}, {"???", cpu.nop, cpu.imp, 3}, {"EOR", cpu.eor, cpu.zp0, 3}, {"LSR", cpu.lsr, cpu.zp0, 5}, {"???", cpu.xxx, cpu.imp, 5}, {"PHA", cpu.pha, cpu.imp, 3}, {"EOR", cpu.eor, cpu.imm, 2}, {"LSR", cpu.lsr, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 2}, {"JMP", cpu.jmp, cpu.abs, 3}, {"EOR", cpu.eor, cpu.abs, 4}, {"LSR", cpu.lsr, cpu.abs, 6}, {"???", cpu.xxx, cpu.imp, 6},
		{"BVC", cpu.bvc, cpu.rel, 2}, {"EOR", cpu.eor, cpu.izy, 5}, {"???", cpu.xxx, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 8}, {"???", cpu.nop, cpu.imp, 4}, {"EOR", cpu.eor, cpu.zpx, 4}, {"LSR", cpu.lsr, cpu.zpx, 6}, {"???", cpu.xxx, cpu.imp, 6}, {"CLI", cpu.cli, cpu.imp, 2}, {"EOR", cpu.eor, cpu.aby, 4}, {"???", cpu.nop, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 7}, {"???", cpu.nop, cpu.imp, 4}, {"EOR", cpu.eor, cpu.abx, 4}, {"LSR", cpu.lsr, cpu.abx, 7}, {"???", cpu.xxx, cpu.imp, 7},
		{"RTS", cpu.rts, cpu.imp, 6}, {"ADC", cpu.adc, cpu.izx, 6}, {"???", cpu.xxx, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 8}, {"???", cpu.nop, cpu.imp, 3}, {"ADC", cpu.adc, cpu.zp0, 3}, {"ROR", cpu.ror, cpu.zp0, 5}, {"???", cpu.xxx, cpu.imp, 5}, {"PLA", cpu.pla, cpu.imp, 4}, {"ADC", cpu.adc, cpu.imm, 2}, {"ROR", cpu.ror, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 2}, {"JMP", cpu.jmp, cpu.ind, 5}, {"ADC", cpu.adc, cpu.abs, 4}, {"ROR", cpu.ror, cpu.abs, 6}, {"???", cpu.xxx, cpu.imp, 6},
		{"BVS", cpu.bvs, cpu.rel, 2}, {"ADC", cpu.adc, cpu.izy, 5}, {"???", cpu.xxx, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 8}, {"???", cpu.nop, cpu.imp, 4}, {"ADC", cpu.adc, cpu.zpx, 4}, {"ROR", cpu.ror, cpu.zpx, 6}, {"???", cpu.xxx, cpu.imp, 6}, {"SEI", cpu.sei, cpu.imp, 2}, {"ADC", cpu.adc, cpu.aby, 4}, {"???", cpu.nop, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 7}, {"???", cpu.nop, cpu.imp, 4}, {"ADC", cpu.adc, cpu.abx, 4}, {"ROR", cpu.ror, cpu.abx, 7}, {"???", cpu.xxx, cpu.imp, 7},
		{"???", cpu.nop, cpu.imp, 2}, {"STA", cpu.sta, cpu.izx, 6}, {"???", cpu.nop, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 6}, {"STY", cpu.sty, cpu.zp0, 3}, {"STA", cpu.sta, cpu.zp0, 3}, {"STX", cpu.stx, cpu.zp0, 3}, {"???", cpu.xxx, cpu.imp, 3}, {"DEY", cpu.dey, cpu.imp, 2}, {"???", cpu.nop, cpu.imp, 2}, {"TXA", cpu.txa, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 2}, {"STY", cpu.sty, cpu.abs, 4}, {"STA", cpu.sta, cpu.abs, 4}, {"STX", cpu.stx, cpu.abs, 4}, {"???", cpu.xxx, cpu.imp, 4},
		{"BCC", cpu.bcc, cpu.rel, 2}, {"STA", cpu.sta, cpu.izy, 6}, {"???", cpu.xxx, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 6}, {"STY", cpu.sty, cpu.zpx, 4}, {"STA", cpu.sta, cpu.zpx, 4}, {"STX", cpu.stx, cpu.zpy, 4}, {"???", cpu.xxx, cpu.imp, 4}, {"TYA", cpu.tya, cpu.imp, 2}, {"STA", cpu.sta, cpu.aby, 5}, {"TXS", cpu.txs, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 5}, {"???", cpu.nop, cpu.imp, 5}, {"STA", cpu.sta, cpu.abx, 5}, {"???", cpu.xxx, cpu.imp, 5}, {"???", cpu.xxx, cpu.imp, 5},
		{"LDY", cpu.ldy, cpu.imm, 2}, {"LDA", cpu.lda, cpu.izx, 6}, {"LDX", cpu.ldx, cpu.imm, 2}, {"???", cpu.xxx, cpu.imp, 6}, {"LDY", cpu.ldy, cpu.zp0, 3}, {"LDA", cpu.lda, cpu.zp0, 3}, {"LDX", cpu.ldx, cpu.zp0, 3}, {"???", cpu.xxx, cpu.imp, 3}, {"TAY", cpu.tay, cpu.imp, 2}, {"LDA", cpu.lda, cpu.imm, 2}, {"TAX", cpu.tax, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 2}, {"LDY", cpu.ldy, cpu.abs, 4}, {"LDA", cpu.lda, cpu.abs, 4}, {"LDX", cpu.ldx, cpu.abs, 4}, {"???", cpu.xxx, cpu.imp, 4},
		{"BCS", cpu.bcs, cpu.rel, 2}, {"LDA", cpu.lda, cpu.izy, 5}, {"???", cpu.xxx, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 5}, {"LDY", cpu.ldy, cpu.zpx, 4}, {"LDA", cpu.lda, cpu.zpx, 4}, {"LDX", cpu.ldx, cpu.zpy, 4}, {"???", cpu.xxx, cpu.imp, 4}, {"CLV", cpu.clv, cpu.imp, 2}, {"LDA", cpu.lda, cpu.aby, 4}, {"TSX", cpu.tsx, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 4}, {"LDY", cpu.ldy, cpu.abx, 4}, {"LDA", cpu.lda, cpu.abx, 4}, {"LDX", cpu.ldx, cpu.aby, 4}, {"???", cpu.xxx, cpu.imp, 4},
		{"CPY", cpu.cpy, cpu.imm, 2}, {"CMP", cpu.cmp, cpu.izx, 6}, {"???", cpu.nop, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 8}, {"CPY", cpu.cpy, cpu.zp0, 3}, {"CMP", cpu.cmp, cpu.zp0, 3}, {"DEC", cpu.dec, cpu.zp0, 5}, {"???", cpu.xxx, cpu.imp, 5}, {"INY", cpu.iny, cpu.imp, 2}, {"CMP", cpu.cmp, cpu.imm, 2}, {"DEX", cpu.dex, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 2}, {"CPY", cpu.cpy, cpu.abs, 4}, {"CMP", cpu.cmp, cpu.abs, 4}, {"DEC", cpu.dec, cpu.abs, 6}, {"???", cpu.xxx, cpu.imp, 6},
		{"BNE", cpu.bne, cpu.rel, 2}, {"CMP", cpu.cmp, cpu.izy, 5}, {"???", cpu.xxx, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 8}, {"???", cpu.nop, cpu.imp, 4}, {"CMP", cpu.cmp, cpu.zpx, 4}, {"DEC", cpu.dec, cpu.zpx, 6}, {"???", cpu.xxx, cpu.imp, 6}, {"CLD", cpu.cld, cpu.imp, 2}, {"CMP", cpu.cmp, cpu.aby, 4}, {"NOP", cpu.nop, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 7}, {"???", cpu.nop, cpu.imp, 4}, {"CMP", cpu.cmp, cpu.abx, 4}, {"DEC", cpu.dec, cpu.abx, 7}, {"???", cpu.xxx, cpu.imp, 7},
		{"CPX", cpu.cpx, cpu.imm, 2}, {"SBC", cpu.sbc, cpu.izx, 6}, {"???", cpu.nop, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 8}, {"CPX", cpu.cpx, cpu.zp0, 3}, {"SBC", cpu.sbc, cpu.zp0, 3}, {"INC", cpu.inc, cpu.zp0, 5}, {"???", cpu.xxx, cpu.imp, 5}, {"INX", cpu.inx, cpu.imp, 2}, {"SBC", cpu.sbc, cpu.imm, 2}, {"NOP", cpu.nop, cpu.imp, 2}, {"???", cpu.sbc, cpu.imp, 2}, {"CPX", cpu.cpx, cpu.abs, 4}, {"SBC", cpu.sbc, cpu.abs, 4}, {"INC", cpu.inc, cpu.abs, 6}, {"???", cpu.xxx, cpu.imp, 6},
		{"BEQ", cpu.beq, cpu.rel, 2}, {"SBC", cpu.sbc, cpu.izy, 5}, {"???", cpu.xxx, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 8}, {"???", cpu.nop, cpu.imp, 4}, {"SBC", cpu.sbc, cpu.zpx, 4}, {"INC", cpu.inc, cpu.zpx, 6}, {"???", cpu.xxx, cpu.imp, 6}, {"SED", cpu.sed, cpu.imp, 2}, {"SBC", cpu.sbc, cpu.aby, 4}, {"NOP", cpu.nop, cpu.imp, 2}, {"???", cpu.xxx, cpu.imp, 7}, {"???", cpu.nop, cpu.imp, 4}, {"SBC", cpu.sbc, cpu.abx, 4}, {"INC", cpu.inc, cpu.abx, 7}, {"???", cpu.xxx, cpu.imp, 7},
	}
}
