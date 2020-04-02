/*
 * locale.go - locale
 *
 * godec128 - go dec128 (for 12-bit decimal fixed point) library
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
    "bytes"
    "strings"
    "strconv"
    "unicode/utf8"
    "github.com/matszpk/goint128"
)

// format 128-bit decimal fixed point including locale
func (a UDec128) LocaleFormatBytes(lang string, tenPow uint,
                                trimZeroes, noSep1000 bool) []byte {
    l := goint128.GetLocFmt(lang)
    s := a.FormatBytes(tenPow, trimZeroes)
    var os bytes.Buffer
    slen := len(s)
    os.Grow(slen*3)
    commaIdx := bytes.LastIndexByte(s, '.')
    if commaIdx==-1 {
        commaIdx = slen
    }
    ti := commaIdx
    i := commaIdx
    if !l.Sep100and1000 {
        ti = (commaIdx)%3
        if ti==0 { ti=3 }
    }
    for k:=0; k < commaIdx; k++ {
        r := s[k]
        if r>='0' && r<='9' {
            os.WriteRune(l.Digits[r-'0'])
        }
        if !noSep1000 && i!=1 {
            if !l.Sep100and1000 || ti<=3 {
                ti--
                if ti==0 {
                    os.WriteRune(l.Sep1000)
                    ti = 3
                }
            } else {
                ti--
                if (ti-3)&1==0 {
                    os.WriteRune(l.Sep1000)
                }
            }
        }
        i--
    }
    // comma
    if commaIdx!=slen {
        os.WriteRune(l.Comma)
        for i = commaIdx+1; i < slen; i++ {
            os.WriteRune(l.Digits[s[i]-'0'])
        }
    }
    return os.Bytes()
}

// format 128-bit decimal fixed point including locale
func (a UDec128) LocaleFormat(lang string, tenPow uint, trimZeroes, noSep1000 bool) string {
    l := goint128.GetLocFmt(lang)
    s := a.FormatBytes(tenPow, trimZeroes)
    var os strings.Builder
    slen := len(s)
    os.Grow(slen*3)
    commaIdx := bytes.LastIndexByte(s, '.')
    if commaIdx==-1 {
        commaIdx = slen
    }
    ti := commaIdx
    i := commaIdx
    if !l.Sep100and1000 {
        ti = (commaIdx)%3
        if ti==0 { ti=3 }
    }
    for k:=0; k < commaIdx; k++ {
        r := s[k]
        if r>='0' && r<='9' {
            os.WriteRune(l.Digits[r-'0'])
        }
        if !noSep1000 && i!=1 {
            if !l.Sep100and1000 || ti<=3 {
                ti--
                if ti==0 {
                    os.WriteRune(l.Sep1000)
                    ti = 3
                }
            } else {
                ti--
                if (ti-3)&1==0 {
                    os.WriteRune(l.Sep1000)
                }
            }
        }
        i--
    }
    // comma
    if commaIdx!=slen {
        os.WriteRune(l.Comma)
        for i = commaIdx+1; i < slen; i++ {
            os.WriteRune(l.Digits[s[i]-'0'])
        }
    }
    return os.String()
}

// parse decimal fixed point from string and return value and error (nil if no error)
func LocaleParseUDec128(lang, str string, tenPow uint, rounding bool) (UDec128, error) {
    l := goint128.GetLocFmt(lang)
    if len(str)==0 { return UDec128{}, strconv.ErrSyntax }
    
    os := make([]byte, 0, len(str))
    for _, r := range str {
        if r>='0' && r<='9' {
            // if standard digits
            os = append(os, byte(r))
        } else if r!=l.Sep1000 && r!=l.Sep1000_2 && r!=l.Comma {
            // if non-standard digit
            dig:=0
            found := false
            for ; dig<=9; dig++ {
                if l.Digits[dig]==r {
                    found = true
                    break
                }
            }
            if !found { return UDec128{}, strconv.ErrSyntax }
            os = append(os, '0'+byte(dig))
        } else if r==l.Comma {
            os = append(os, '.')
        }
        // otherwise skip sep1000
    }
    return ParseUDec128Bytes(os, tenPow, rounding)
}

// parse decimal fixed point from string and return value and error (nil if no error)
func LocaleParseUDec128Bytes(lang string, strInput []byte,
                             tenPow uint, rounding bool) (UDec128, error) {
    l := goint128.GetLocFmt(lang)
    if len(strInput)==0 { return UDec128{}, strconv.ErrSyntax }
    
    os := make([]byte, 0, len(strInput))
    str := strInput
    for len(str)>0 {
        r, size := utf8.DecodeRune(str)
        if r>='0' && r<='9' {
            // if standard digits
            os = append(os, byte(r))
        } else if r!=l.Sep1000 && r!=l.Sep1000_2 && r!=l.Comma {
            // if non-standard digit
            dig:=0
            found := false
            for ; dig<=9; dig++ {
                if l.Digits[dig]==r {
                    found = true
                    break
                }
            }
            if !found { return UDec128{}, strconv.ErrSyntax }
            os = append(os, '0'+byte(dig))
        } else if r==l.Comma {
            os = append(os, '.')
        }
        // otherwise skip sep1000
        str = str[size:]
    }
    return ParseUDec128Bytes(os, tenPow, rounding)
}

