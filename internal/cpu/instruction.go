package cpu

const (
	// address modes
	imp = "IMP"; imm = "IMM"
	zp0 = "ZP0"; zpx = "ZPX"; zpy = "ZPY"
	rel = "REL"
	abs = "ABS"; abx = "ABX"; aby = "ABY"
	ind = "IND"; izx = "IZX"; izy = "IZY"
	
	// operations
	adc = "ADC"; and = "AND"; asl = "ASL"
	bcc = "BCC"; bcs = "BCS"; beq = "BEQ"; bit = "BIT"; bmi = "BMI"; bne = "BNE"; bpl = "BPL"; brk = "BRK";  bvc = "BVC"; bvs = "BVS"
	clc = "CLC"; cld = "CLD"; cli = "CLI"; clv = "CLV"; cmp = "CMP"; cpx = "CPX"; cpy = "CPY"
	dec = "DEC"; dex = "DEX"; dey = "DEY"
	eor = "EOR"
	inc = "INC"; inx = "INX"; iny = "INY"
	jmp = "JMP"; jsr = "JSR"
	lda = "LDA"; ldx = "LDX"; ldy = "LDY"; lsr = "LSR"
	nop = "NOP"
	ora = "ORA"
	pha = "PHA"; php = "PHP"; pla = "PLA"; plp = "PLP"
	rol = "ROL"; ror = "ROR"; rti = "RTI"; rts = "RTS"
	sbc = "SBC"; sec = "SEC"; sed = "SED"; sei = "SEI"; sta = "STA"; stx = "STX"; sty = "STY"
	tax = "TAX"; tay = "TAY"; tsx = "TSX"; txa = "TXA"; txs = "TXS"; tya = "TYA"
	xxx = "???"
)

type instruction struct {
	operation      string
	addressMode    string
	performOp      func() uint8
	setAddressMode func() uint8
	cycles         uint8
}

type mos6502LookupTable [256]instruction

func buildMos502LookupTable(c *Mos6502) mos6502LookupTable {
	return mos6502LookupTable{
		{brk, imm, c.brk, c.imm, 7}, {ora, izx, c.ora, c.izx, 6}, {xxx, imp, c.xxx, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 8}, {xxx, imp, c.nop, c.imp, 3}, {ora, zp0, c.ora, c.zp0, 3}, {asl, zp0, c.asl, c.zp0, 5}, {xxx, imp, c.xxx, c.imp, 5}, {php, imp, c.php, c.imp, 3}, {ora, imm, c.ora, c.imm, 2}, {asl, imp, c.asl, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 2}, {xxx, imp, c.nop, c.imp, 4}, {ora, abs, c.ora, c.abs, 4}, {asl, abs, c.asl, c.abs, 6}, {xxx, imp, c.xxx, c.imp, 6},
		{bpl, rel, c.bpl, c.rel, 2}, {ora, izy, c.ora, c.izy, 5}, {xxx, imp, c.xxx, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 8}, {xxx, imp, c.nop, c.imp, 4}, {ora, zpx, c.ora, c.zpx, 4}, {asl, zpx, c.asl, c.zpx, 6}, {xxx, imp, c.xxx, c.imp, 6}, {clc, imp, c.clc, c.imp, 2}, {ora, aby, c.ora, c.aby, 4}, {xxx, imp, c.nop, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 7}, {xxx, imp, c.nop, c.imp, 4}, {ora, abx, c.ora, c.abx, 4}, {asl, abx, c.asl, c.abx, 7}, {xxx, imp, c.xxx, c.imp, 7},
		{jsr, abs, c.jsr, c.abs, 6}, {and, izx, c.and, c.izx, 6}, {xxx, imp, c.xxx, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 8}, {bit, zp0, c.bit, c.zp0, 3}, {and, zp0, c.and, c.zp0, 3}, {rol, zp0, c.rol, c.zp0, 5}, {xxx, imp, c.xxx, c.imp, 5}, {plp, imp, c.plp, c.imp, 4}, {and, imm, c.and, c.imm, 2}, {rol, imp, c.rol, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 2}, {bit, abs, c.bit, c.abs, 4}, {and, abs, c.and, c.abs, 4}, {rol, abs, c.rol, c.abs, 6}, {xxx, imp, c.xxx, c.imp, 6},
		{bmi, rel, c.bmi, c.rel, 2}, {and, izy, c.and, c.izy, 5}, {xxx, imp, c.xxx, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 8}, {xxx, imp, c.nop, c.imp, 4}, {and, zpx, c.and, c.zpx, 4}, {rol, zpx, c.rol, c.zpx, 6}, {xxx, imp, c.xxx, c.imp, 6}, {sec, imp, c.sec, c.imp, 2}, {and, aby, c.and, c.aby, 4}, {xxx, imp, c.nop, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 7}, {xxx, imp, c.nop, c.imp, 4}, {and, abx, c.and, c.abx, 4}, {rol, abx, c.rol, c.abx, 7}, {xxx, imp, c.xxx, c.imp, 7},
		{rti, imp, c.rti, c.imp, 6}, {eor, izx, c.eor, c.izx, 6}, {xxx, imp, c.xxx, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 8}, {xxx, imp, c.nop, c.imp, 3}, {eor, zp0, c.eor, c.zp0, 3}, {lsr, zp0, c.lsr, c.zp0, 5}, {xxx, imp, c.xxx, c.imp, 5}, {pha, imp, c.pha, c.imp, 3}, {eor, imm, c.eor, c.imm, 2}, {lsr, imp, c.lsr, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 2}, {jmp, abs, c.jmp, c.abs, 3}, {eor, abs, c.eor, c.abs, 4}, {lsr, abs, c.lsr, c.abs, 6}, {xxx, imp, c.xxx, c.imp, 6},
		{bvc, rel, c.bvc, c.rel, 2}, {eor, izy, c.eor, c.izy, 5}, {xxx, imp, c.xxx, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 8}, {xxx, imp, c.nop, c.imp, 4}, {eor, zpx, c.eor, c.zpx, 4}, {lsr, zpx, c.lsr, c.zpx, 6}, {xxx, imp, c.xxx, c.imp, 6}, {cli, imp, c.cli, c.imp, 2}, {eor, aby, c.eor, c.aby, 4}, {xxx, imp, c.nop, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 7}, {xxx, imp, c.nop, c.imp, 4}, {eor, abx, c.eor, c.abx, 4}, {lsr, abx, c.lsr, c.abx, 7}, {xxx, imp, c.xxx, c.imp, 7},
		{rts, imp, c.rts, c.imp, 6}, {adc, izx, c.adc, c.izx, 6}, {xxx, imp, c.xxx, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 8}, {xxx, imp, c.nop, c.imp, 3}, {adc, zp0, c.adc, c.zp0, 3}, {ror, zp0, c.ror, c.zp0, 5}, {xxx, imp, c.xxx, c.imp, 5}, {pla, imp, c.pla, c.imp, 4}, {adc, imm, c.adc, c.imm, 2}, {ror, imp, c.ror, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 2}, {jmp, ind, c.jmp, c.ind, 5}, {adc, abs, c.adc, c.abs, 4}, {ror, abs, c.ror, c.abs, 6}, {xxx, imp, c.xxx, c.imp, 6},
		{bvs, rel, c.bvs, c.rel, 2}, {adc, izy, c.adc, c.izy, 5}, {xxx, imp, c.xxx, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 8}, {xxx, imp, c.nop, c.imp, 4}, {adc, zpx, c.adc, c.zpx, 4}, {ror, zpx, c.ror, c.zpx, 6}, {xxx, imp, c.xxx, c.imp, 6}, {sei, imp, c.sei, c.imp, 2}, {adc, aby, c.adc, c.aby, 4}, {xxx, imp, c.nop, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 7}, {xxx, imp, c.nop, c.imp, 4}, {adc, abx, c.adc, c.abx, 4}, {ror, abx, c.ror, c.abx, 7}, {xxx, imp, c.xxx, c.imp, 7},
		{xxx, imp, c.nop, c.imp, 2}, {sta, izx, c.sta, c.izx, 6}, {xxx, imp, c.nop, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 6}, {sty, zp0, c.sty, c.zp0, 3}, {sta, zp0, c.sta, c.zp0, 3}, {stx, zp0, c.stx, c.zp0, 3}, {xxx, imp, c.xxx, c.imp, 3}, {dey, imp, c.dey, c.imp, 2}, {xxx, imp, c.nop, c.imp, 2}, {txa, imp, c.txa, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 2}, {sty, abs, c.sty, c.abs, 4}, {sta, abs, c.sta, c.abs, 4}, {stx, abs, c.stx, c.abs, 4}, {xxx, imp, c.xxx, c.imp, 4},
		{bcc, rel, c.bcc, c.rel, 2}, {sta, izy, c.sta, c.izy, 6}, {xxx, imp, c.xxx, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 6}, {sty, zpx, c.sty, c.zpx, 4}, {sta, zpx, c.sta, c.zpx, 4}, {stx, zpy, c.stx, c.zpy, 4}, {xxx, imp, c.xxx, c.imp, 4}, {tya, imp, c.tya, c.imp, 2}, {sta, aby, c.sta, c.aby, 5}, {txs, imp, c.txs, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 5}, {xxx, imp, c.nop, c.imp, 5}, {sta, abx, c.sta, c.abx, 5}, {xxx, imp, c.xxx, c.imp, 5}, {xxx, imp, c.xxx, c.imp, 5},
		{ldy, imm, c.ldy, c.imm, 2}, {lda, izx, c.lda, c.izx, 6}, {ldx, imm, c.ldx, c.imm, 2}, {xxx, imp, c.xxx, c.imp, 6}, {ldy, zp0, c.ldy, c.zp0, 3}, {lda, zp0, c.lda, c.zp0, 3}, {ldx, zp0, c.ldx, c.zp0, 3}, {xxx, imp, c.xxx, c.imp, 3}, {tay, imp, c.tay, c.imp, 2}, {lda, imm, c.lda, c.imm, 2}, {tax, imp, c.tax, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 2}, {ldy, abs, c.ldy, c.abs, 4}, {lda, abs, c.lda, c.abs, 4}, {ldx, abs, c.ldx, c.abs, 4}, {xxx, imp, c.xxx, c.imp, 4},
		{bcs, rel, c.bcs, c.rel, 2}, {lda, izy, c.lda, c.izy, 5}, {xxx, imp, c.xxx, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 5}, {ldy, zpx, c.ldy, c.zpx, 4}, {lda, zpx, c.lda, c.zpx, 4}, {ldx, zpy, c.ldx, c.zpy, 4}, {xxx, imp, c.xxx, c.imp, 4}, {clv, imp, c.clv, c.imp, 2}, {lda, aby, c.lda, c.aby, 4}, {tsx, imp, c.tsx, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 4}, {ldy, abx, c.ldy, c.abx, 4}, {lda, abx, c.lda, c.abx, 4}, {ldx, aby, c.ldx, c.aby, 4}, {xxx, imp, c.xxx, c.imp, 4},
		{cpy, imm, c.cpy, c.imm, 2}, {cmp, izx, c.cmp, c.izx, 6}, {xxx, imp, c.nop, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 8}, {cpy, zp0, c.cpy, c.zp0, 3}, {cmp, zp0, c.cmp, c.zp0, 3}, {dec, zp0, c.dec, c.zp0, 5}, {xxx, imp, c.xxx, c.imp, 5}, {iny, imp, c.iny, c.imp, 2}, {cmp, imm, c.cmp, c.imm, 2}, {dex, imp, c.dex, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 2}, {cpy, abs, c.cpy, c.abs, 4}, {cmp, abs, c.cmp, c.abs, 4}, {dec, abs, c.dec, c.abs, 6}, {xxx, imp, c.xxx, c.imp, 6},
		{bne, rel, c.bne, c.rel, 2}, {cmp, izy, c.cmp, c.izy, 5}, {xxx, imp, c.xxx, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 8}, {xxx, imp, c.nop, c.imp, 4}, {cmp, zpx, c.cmp, c.zpx, 4}, {dec, zpx, c.dec, c.zpx, 6}, {xxx, imp, c.xxx, c.imp, 6}, {cld, imp, c.cld, c.imp, 2}, {cmp, aby, c.cmp, c.aby, 4}, {nop, imp, c.nop, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 7}, {xxx, imp, c.nop, c.imp, 4}, {cmp, abx, c.cmp, c.abx, 4}, {dec, abx, c.dec, c.abx, 7}, {xxx, imp, c.xxx, c.imp, 7},
		{cpx, imm, c.cpx, c.imm, 2}, {sbc, izx, c.sbc, c.izx, 6}, {xxx, imp, c.nop, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 8}, {cpx, zp0, c.cpx, c.zp0, 3}, {sbc, zp0, c.sbc, c.zp0, 3}, {inc, zp0, c.inc, c.zp0, 5}, {xxx, imp, c.xxx, c.imp, 5}, {inx, imp, c.inx, c.imp, 2}, {sbc, imm, c.sbc, c.imm, 2}, {nop, imp, c.nop, c.imp, 2}, {xxx, imp, c.sbc, c.imp, 2}, {cpx, abs, c.cpx, c.abs, 4}, {sbc, abs, c.sbc, c.abs, 4}, {inc, abs, c.inc, c.abs, 6}, {xxx, imp, c.xxx, c.imp, 6},
		{beq, rel, c.beq, c.rel, 2}, {sbc, izy, c.sbc, c.izy, 5}, {xxx, imp, c.xxx, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 8}, {xxx, imp, c.nop, c.imp, 4}, {sbc, zpx, c.sbc, c.zpx, 4}, {inc, zpx, c.inc, c.zpx, 6}, {xxx, imp, c.xxx, c.imp, 6}, {sed, imp, c.sed, c.imp, 2}, {sbc, aby, c.sbc, c.aby, 4}, {nop, imp, c.nop, c.imp, 2}, {xxx, imp, c.xxx, c.imp, 7}, {xxx, imp, c.nop, c.imp, 4}, {sbc, abx, c.sbc, c.abx, 4}, {inc, abx, c.inc, c.abx, 7}, {xxx, imp, c.xxx, c.imp, 7},
	}
}
