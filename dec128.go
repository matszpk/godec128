/*
 * dec128.go - main fixed decimal int128 routines
 *
 * goint128 - go dec128 (for 12-bit decimal fixed point) library
 * Copyright (C) 2019  Mateusz Szpakowski
 *
 * This library is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public
 * License as published by the Free Software Foundation; either
 * version 2.1 of the License, or (at your option) any later version.
 *
 * This library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public
 * License along with this library; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301  USA
 */

// Package to operate on 128-bit decimal fixed point
package godec128

import (
    "math/bits"
    "strconv"
    "strings"
    "github.com/matszpk/goint128"
)

type UDec128 goint128.UInt128

// add 128-bit unsigned integers
func (a UDec128) Add(b UDec128) UDec128 {
    return UDec128(goint128.UInt128(a).Add(goint128.UInt128(b)))
}

// add 128-bit unsigned integers with carry and return sum and output carry
func (a UDec128) AddC(b UDec128, oldCarry uint64) (UDec128, uint64) {
    v, c := goint128.UInt128(a).AddC(goint128.UInt128(b), oldCarry)
    return UDec128(v), c
}

// add 128-bit unsigned integer and 64-bit unsigned integer
func (a UDec128) Add64(b uint64) UDec128 {
    return UDec128(goint128.UInt128(a).Add64(b))
}

// subtract 128-bit unsigned integers
func (a UDec128) Sub(b UDec128) UDec128 {
    return UDec128(goint128.UInt128(a).Sub(goint128.UInt128(b)))
}

// subtract 128-bit unsigned integers with borrow and return difference and borrow
func (a UDec128) SubB(b UDec128, oldBorrow uint64) (UDec128, uint64) {
    v, br := goint128.UInt128(a).SubB(goint128.UInt128(b), oldBorrow)
    return UDec128(v), br
}

// subtract 64-bit unsigned from 128-bit unsigned integer
func (a UDec128) Sub64(b uint64) UDec128 {
    return UDec128(goint128.UInt128(a).Sub64(b))
}

func (a UDec128) Cmp(b UDec128) int {
    return goint128.UInt128(a).Cmp(goint128.UInt128(b))
}

var uint64_powers []uint64 = []uint64{
    1,
    10,
    100,
    1000,
    10000,
    100000,
    1000000,
    10000000,
    100000000,
    1000000000,
    10000000000,
    100000000000,
    1000000000000,
    10000000000000,
    100000000000000,
    1000000000000000,
}

func uint128_64DivFullR(hi, lo goint128.UInt128, b uint64,
                        rounding bool) goint128.UInt128 {
    if b==1 { return lo }
    var borrow uint64
    lza := 0
    if hi[0]==0 && hi[1]==0 {
        lza = 128
    } else if hi[1]!=0 {
        lza = bits.LeadingZeros64(hi[1])
    } else {
        lza = bits.LeadingZeros64(hi[0])+64
    }
    lzb := bits.LeadingZeros64(b)+64
    sh := uint(lza-lzb)
    pos := int(128-sh)
    // shift A (lo,hi) by shift (move to highest bit of b)
    var tlo, thi goint128.UInt128
    if sh!=128 {
        tlo = lo.Shl(sh)
        thi = hi.Shl(sh)
        if sh!=0 {
            tmp := lo.Shr(128-sh)
            thi[0] |= tmp[0]
            thi[1] |= tmp[1]
        }
    } else {
        thi = lo
        tlo[0], tlo[1] = 0, 0
    }
    // main loop
    var tmp goint128.UInt128
    c := goint128.UInt128{0,0}
    for ; pos>0; pos-- {
        tmp[0], borrow = goint128.Sub64(thi[0], b, 0)
        tmp[1], borrow = goint128.Sub64(thi[1], 0, borrow)
        c[1] = (c[0]>>63) | (c[1]<<1) // shift
        c[0] <<= 1
        if borrow==0 {
            thi = tmp
            c[0] |= 1
        }
        // shift T (shifted copy of A)
        thi[1] = (thi[0]>>63) | (thi[1]<<1) // shift
        thi[0] = (tlo[1]>>63) | (thi[0]<<1)
        tlo[1] = (tlo[0]>>63) | (tlo[1]<<1)
        tlo[0] <<= 1
    }
    // last iteration
    tmp[0], borrow = goint128.Sub64(thi[0], b, 0)
    tmp[1], borrow = goint128.Sub64(thi[1], 0, borrow)
    c[1] = (c[0]>>63) | (c[1]<<1) // shift
    c[0] <<= 1
    if borrow==0 {
        thi = tmp
        c[0] |= 1
    }
    if rounding && thi[0]>=(b>>1) { // rounding
        c = c.Add64(1)
    }
    return c
}

func (a UDec128) Mul(b UDec128, tenPow uint, rounding bool) UDec128 {
    chi, clo := goint128.UInt128(a).MulFull(goint128.UInt128(b))
    // divide by ten power
    return UDec128(uint128_64DivFullR(chi, clo, uint64_powers[tenPow], rounding))
}

func (a UDec128) Mul64(b uint64) UDec128 {
    return UDec128(goint128.UInt128(a).Mul64(b))
}

// fixed point is in 10**(tenPow*2)
func (a UDec128) MulFull(b UDec128) (UDec128, UDec128) {
    chi, clo := goint128.UInt128(a).MulFull(goint128.UInt128(b))
    return UDec128(chi), UDec128(clo)
}

// shift 128-bit unsigned integer left by b bits
func (a UDec128) Shl(b uint) UDec128 {
    return UDec128(goint128.UInt128(a).Shl(b))
}

func (a UDec128) Shr(b uint) UDec128 {
    return UDec128(goint128.UInt128(a).Shr(b))
}

func (a UDec128) Div(b UDec128, tenPow uint) (UDec128, UDec128) {
    // multiply by tenPowers
    chi, clo := goint128.UInt128(a).MulFull(goint128.UInt128{uint64_powers[tenPow], 0})
    q, r := goint128.UInt128DivFull(chi, clo, goint128.UInt128(b))
    return UDec128(q), UDec128(r)
}

func (a UDec128) Div64(b uint64) (UDec128, uint64) {
    q, r := goint128.UInt128(a).Div64(b)
    return UDec128(q), r
}

// fixed point is in 10**(tenPow*2)
func UDec128DivFull(hi, lo, b UDec128) (UDec128, UDec128) {
    q, r := goint128.UInt128DivFull(goint128.UInt128(hi), goint128.UInt128(lo),
                                    goint128.UInt128(b))
    return UDec128(q), UDec128(r)
}

var zeroPart string = "0.000000000000000000000000000"

func (a UDec128) Format(tenPow uint, trimZeroes bool) string {
    str := goint128.UInt128(a).Format()
    if tenPow==0 { return str }
    slen := len(str)
    i := slen-1
    if slen <= int(tenPow) {
        if trimZeroes {
            for ; i>=0; i-- {
                if str[i]!='0' { break }
            }
        }
        return zeroPart[:2+int(tenPow)-slen] + str[:i]
    }
    if trimZeroes {
        for ; i>=int(tenPow); i-- {
            if str[i]!='0' { break }
        }
    }
    return str[:slen-int(tenPow)]+"."+str[slen-int(tenPow):i]
}

func ParseUDec128(str string, tenPow uint) (UDec128, error) {
    if tenPow==0 {
        v, err := goint128.ParseUInt128(str)
        return UDec128(v), err
    }
    slen := len(str)
    commaIdx := strings.LastIndexByte(str, '.')
    if commaIdx==-1 {
        // comma not found
        v, err := goint128.ParseUInt128(str)
        if err!=nil { return UDec128(v), err }
        chi, clo := v.MulFull(goint128.UInt128{uint64_powers[tenPow], 0})
        if chi[0]!=0 || chi[1]!=0 {
            return UDec128{}, strconv.ErrRange
        }
        return UDec128(clo), nil
    }
    if slen-(commaIdx+1) >= int(tenPow) {
        //  more than in fraction
        realSlen := commaIdx+1+int(tenPow)
        s2 := str[:commaIdx] + str[commaIdx+1:realSlen]
        v, err := goint128.ParseUInt128(s2)
        if err!=nil { return UDec128{}, err }
        // rounding
        if realSlen!=slen && str[realSlen]>='5' {
            v = v.Add64(1) // add rounding
        }
        // check last part of string
        for i:=realSlen; i<slen; i++ {
            if str[i]<'0' || str[i]>'9' {
                return UDec128{}, strconv.ErrSyntax
            }
        }
        return UDec128(v), nil
    } else {
        // less than in fraction
        s2 := str[:commaIdx] + str[commaIdx+1:]
        v, err := goint128.ParseUInt128(s2)
        if err!=nil { return UDec128{}, err }
        pow10ForVal := int(tenPow) - (slen-(commaIdx+1))
        chi, clo := v.MulFull(goint128.UInt128{uint64_powers[pow10ForVal], 0})
        if chi[0]!=0 || chi[1]!=0 {
            return UDec128{}, strconv.ErrRange
        }
        return UDec128(clo), nil
    }
    return UDec128{}, nil
}

var float64_revpowers []float64 = []float64{
    1,
    0.1,
    0.01,
    0.001,
    0.0001,
    0.00001,
    0.000001,
    0.0000001,
    0.00000001,
    0.000000001,
    0.0000000001,
    0.00000000001,
    0.000000000001,
    0.0000000000001,
    0.00000000000001,
    0.000000000000001,
}

func (a UDec128) ToFloat64(tenPow int) float64 {
    return goint128.UInt128(a).ToFloat64()*float64_revpowers[tenPow]
}

func Float64ToUDec128(a float64, tenPow int) (UDec128, error) {
    r, err := goint128.Float64ToUInt128(a*float64(uint64_powers[tenPow]))
    return UDec128(r), err
}
