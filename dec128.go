/*
 * dec128.go - main fixed decimal int128 routines
 *
 * godec128 - go dec128 (for 12-bit decimal fixed point) library
 * Copyright (C) 2020  Mateusz Szpakowski
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
    "bytes"
    "math/bits"
    "strconv"
    "strings"
    "github.com/matszpk/goint128"
)

type UDec128 goint128.UInt128

// add 128-bit decimal fixed points
func (a UDec128) Add(b UDec128) UDec128 {
    return UDec128(goint128.UInt128(a).Add(goint128.UInt128(b)))
}

// add 128-bit decimal fixed points with carry and return sum and output carry
func (a UDec128) AddC(b UDec128, oldCarry uint64) (UDec128, uint64) {
    v, c := goint128.UInt128(a).AddC(goint128.UInt128(b), oldCarry)
    return UDec128(v), c
}

// subtract 128-bit decimal fixed points
func (a UDec128) Sub(b UDec128) UDec128 {
    return UDec128(goint128.UInt128(a).Sub(goint128.UInt128(b)))
}

// subtract 128-bit decimal fixed points with borrow and return difference and borrow
func (a UDec128) SubB(b UDec128, oldBorrow uint64) (UDec128, uint64) {
    v, br := goint128.UInt128(a).SubB(goint128.UInt128(b), oldBorrow)
    return UDec128(v), br
}

// compare 128-bit decimal fixed points and return 0 if they equal,
// 1 if first is greater than second, or -1 if first is lesser than second
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
    10000000000000000,
    100000000000000000,
    1000000000000000000,
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

// multiply 128-bit decimal fixed points and return lower 128 bits value
func (a UDec128) Mul(b UDec128, precision uint, rounding bool) UDec128 {
    chi, clo := goint128.UInt128(a).MulFull(goint128.UInt128(b))
    // divide by ten power
    return UDec128(uint128_64DivFullR(chi, clo, uint64_powers[precision], rounding))
}

// multiply 128-bit decimal fixed point and 64-bit unsigned integer and
// return lower 128 bits product
func (a UDec128) Mul64(b uint64) UDec128 {
    return UDec128(goint128.UInt128(a).Mul64(b))
}

// multiply 128-bit decimal fixed point and return high and lower product
// integer part is multiplied by 10**precision
func (a UDec128) MulFull(b UDec128) (UDec128, UDec128) {
    chi, clo := goint128.UInt128(a).MulFull(goint128.UInt128(b))
    return UDec128(chi), UDec128(clo)
}

// shift 128-bit decimal fixed point left by b bits
func (a UDec128) Shl(b uint) UDec128 {
    return UDec128(goint128.UInt128(a).Shl(b))
}

// shift 128-bit decimal fixed point right by b bits
func (a UDec128) Shr(b uint) UDec128 {
    return UDec128(goint128.UInt128(a).Shr(b))
}

// divide 128-bit decimal fixed points
func (a UDec128) Div(b UDec128, precision uint) UDec128 {
    // multiply by precisioners
    chi, clo := goint128.UInt128(a).MulFull(goint128.UInt128{uint64_powers[precision], 0})
    q, _ := goint128.UInt128DivFull(chi, clo, goint128.UInt128(b))
    return UDec128(q)
}

// divide 128-bit unsigned integer by 64-bit unsigned integer
func (a UDec128) Div64(b uint64) UDec128 {
    q, _ := goint128.UInt128(a).Div64(b)
    return UDec128(q)
}

// fixed point is in 10**(precision*2)
func UDec128DivFull(hi, lo, b UDec128) UDec128 {
    q, _ := goint128.UInt128DivFull(goint128.UInt128(hi), goint128.UInt128(lo),
                                    goint128.UInt128(b))
    return UDec128(q)
}

var zeroPart []byte = []byte("0.000000000000000000000000000")

func (a UDec128) Format(precision uint, trimZeroes bool) string {
    if a[0]==0 && a[1]==0 { return "0.0" }
    if precision==0 { return goint128.UInt128(a).Format() }
    str := goint128.UInt128(a).FormatBytes()
    slen := len(str)
    i := slen
    if slen <= int(precision) {
        if trimZeroes {
            for i--; i>=0; i-- {
                if str[i]!='0' { break }
            }
            i++
        }
        var os strings.Builder
        os.Write(zeroPart[:2+int(precision)-slen])
        os.Write(str[:i])
        return os.String()
    }
    if trimZeroes {
        for i--; i>=slen-int(precision); i-- {
            if str[i]!='0' { break }
        }
        i++
    }
    var os strings.Builder
    os.Grow(i)
    os.Write(str[:slen-int(precision)])
    os.WriteByte('.')
    os.Write(str[slen-int(precision):i])
    return os.String()
}

func (a UDec128) FormatBytes(precision uint, trimZeroes bool) []byte {
    if a[0]==0 && a[1]==0 { return zeroPart[:3] }
    if precision==0 { return goint128.UInt128(a).FormatBytes() }
    str := goint128.UInt128(a).FormatBytes()
    slen := len(str)
    i := slen
    if slen <= int(precision) {
        if trimZeroes {
            for i--; i>=0; i-- {
                if str[i]!='0' { break }
            }
            i++
        }
        l := 2+int(precision)-slen
        os := make([]byte, l+i)
        copy(os[:l], zeroPart[:l])
        copy(os[l:], str[:i])
        return os
    }
    if trimZeroes {
        for i--; i>=slen-int(precision); i-- {
            if str[i]!='0' { break }
        }
        i++
    }
    os := make([]byte, i+1)
    l := slen-int(precision)
    copy(os[:l], str[:l])
    os[l] = '.'
    copy(os[l+1:], str[slen-int(precision):i])
    return os
}


func ParseUDec128(str string, precision uint, rounding bool) (UDec128, error) {
    if precision==0 {
        v, err := goint128.ParseUInt128(str)
        return UDec128(v), err
    }
    slen := len(str)
    epos := strings.LastIndexByte(str, 'e')
    if epos!=-1 {
        // parse exponent
        if epos+1==slen {
            return UDec128{}, strconv.ErrSyntax
        }
        // sign of exponent
        endOfMantisa := epos
        epos++
        exponent, err := strconv.ParseInt(str[epos:], 10, 8)
        if err!=nil { return UDec128{}, err }
        
        if exponent!=0 {
            mantisa := str[:endOfMantisa]
            commaPos := strings.IndexByte(mantisa, '.')
            // move comma
            if commaPos==-1 { commaPos = endOfMantisa }
            
            newCommaPos := commaPos + int(exponent)
            
            i := 0
            for ; str[i]=='0' || str[i]=='.'; i++ {
                if str[i]=='.' { continue }
                newCommaPos--
            } // skip first zero
            //fmt.Println("NewCommaPos:", newCommaPos, commaPos, exponent, i)
            var sb strings.Builder
            // add zeroes
            if newCommaPos<0 {
                sb.WriteRune('.')
                for ; newCommaPos<0;  newCommaPos++ {
                    sb.WriteRune('0')
                }
            } else if newCommaPos>0 {
                for ; newCommaPos>0;  newCommaPos-- {
                    if str[i]=='.' { i++ }
                    if i<endOfMantisa {
                        sb.WriteByte(str[i])
                        i++
                    } else {
                        sb.WriteRune('0')
                    }
                }
                if i<endOfMantisa {
                    sb.WriteRune('.') // append new comma
                }
            }
            // to end of mantisa
            for ; i<endOfMantisa; i++ {
                if str[i]=='.' { i++ }
                if i<endOfMantisa { sb.WriteByte(str[i]) }
            }
            
            //fmt.Println("new str:", sb.String())
            str = sb.String()
            slen = len(str)
            if slen==0 { return UDec128{}, nil }
        } else {
            str = str[:endOfMantisa]
            slen = len(str)
        }
    }
    
    commaIdx := strings.LastIndexByte(str, '.')
    if commaIdx==-1 {
        // comma not found
        v, err := goint128.ParseUInt128(str)
        if err!=nil { return UDec128(v), err }
        chi, clo := v.MulFull(goint128.UInt128{uint64_powers[precision], 0})
        if chi[0]!=0 || chi[1]!=0 {
            return UDec128{}, strconv.ErrRange
        }
        return UDec128(clo), nil
    }
    if slen-(commaIdx+1) >= int(precision) {
        //  more than in fraction
        realSlen := commaIdx+1+int(precision)
        s2 := str[:commaIdx] + str[commaIdx+1:realSlen]
        v, err := goint128.ParseUInt128(s2)
        if err!=nil { return UDec128{}, err }
        // rounding
        if rounding && realSlen!=slen && str[realSlen]>='5' {
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
        pow10ForVal := int(precision) - (slen-(commaIdx+1))
        chi, clo := v.MulFull(goint128.UInt128{uint64_powers[pow10ForVal], 0})
        if chi[0]!=0 || chi[1]!=0 {
            return UDec128{}, strconv.ErrRange
        }
        return UDec128(clo), nil
    }
    return UDec128{}, nil
}

func ParseUDec128Bytes(str []byte, precision uint, rounding bool) (UDec128, error) {
    if precision==0 {
        v, err := goint128.ParseUInt128Bytes(str)
        return UDec128(v), err
    }
    
    slen := len(str)
    epos := bytes.LastIndexByte(str, 'e')
    if epos!=-1 {
        // parse exponent
        if epos+1==slen {
            return UDec128{}, strconv.ErrSyntax
        }
        // sign of exponent
        endOfMantisa := epos
        epos++
        exponent, err := strconv.ParseInt(string(str[epos:]), 10, 8)
        if err!=nil { return UDec128{}, err }
        
        if exponent!=0 {
            mantisa := str[:endOfMantisa]
            commaPos := bytes.IndexByte(mantisa, '.')
            // move comma
            if commaPos==-1 { commaPos = endOfMantisa }
            
            newCommaPos := commaPos + int(exponent)
            
            i := 0
            for ; str[i]=='0' || str[i]=='.'; i++ {
                if str[i]=='.' { continue }
                newCommaPos--
            } // skip first zero
            //fmt.Println("NewCommaPos:", newCommaPos, commaPos, exponent, i)
            var sb bytes.Buffer
            // add zeroes
            if newCommaPos<0 {
                sb.WriteRune('.')
                for ; newCommaPos<0;  newCommaPos++ {
                    sb.WriteRune('0')
                }
            } else if newCommaPos>0 {
                for ; newCommaPos>0;  newCommaPos-- {
                    if str[i]=='.' { i++ }
                    if i<endOfMantisa {
                        sb.WriteByte(str[i])
                        i++
                    } else {
                        sb.WriteRune('0')
                    }
                }
                if i<endOfMantisa {
                    sb.WriteRune('.') // append new comma
                }
            }
            // to end of mantisa
            for ; i<endOfMantisa; i++ {
                if str[i]=='.' { i++ }
                if i<endOfMantisa { sb.WriteByte(str[i]) }
            }
            
            //fmt.Println("new str:", sb.String())
            str = sb.Bytes()
            slen = len(str)
            if slen==0 { return UDec128{}, nil }
        } else {
            str = str[:endOfMantisa]
            slen = len(str)
        }
    }
    
    commaIdx := bytes.LastIndexByte(str, '.')
    if commaIdx==-1 {
        // comma not found
        v, err := goint128.ParseUInt128Bytes(str)
        if err!=nil { return UDec128(v), err }
        chi, clo := v.MulFull(goint128.UInt128{uint64_powers[precision], 0})
        if chi[0]!=0 || chi[1]!=0 {
            return UDec128{}, strconv.ErrRange
        }
        return UDec128(clo), nil
    }
    if slen-(commaIdx+1) >= int(precision) {
        //  more than in fraction
        realSlen := commaIdx+1+int(precision)
        s2 := make([]byte, realSlen-1)
        copy(s2[:commaIdx], str[:commaIdx])
        copy(s2[commaIdx:], str[commaIdx+1:])
        v, err := goint128.ParseUInt128Bytes(s2)
        if err!=nil { return UDec128{}, err }
        // rounding
        if rounding && realSlen!=slen && str[realSlen]>='5' {
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
        s2 := make([]byte, slen-1)
        copy(s2[:commaIdx], str[:commaIdx])
        copy(s2[commaIdx:], str[commaIdx+1:])
        v, err := goint128.ParseUInt128Bytes(s2)
        if err!=nil { return UDec128{}, err }
        pow10ForVal := int(precision) - (slen-(commaIdx+1))
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
    0.0000000000000001,
    0.00000000000000001,
    0.000000000000000001,
}

func (a UDec128) ToFloat64(precision uint) float64 {
    return goint128.UInt128(a).ToFloat64()*float64_revpowers[precision]
}

func Float64ToUDec128(a float64, precision uint) (UDec128, error) {
    r, err := goint128.Float64ToUInt128(a*float64(uint64_powers[precision]))
    return UDec128(r), err
}
