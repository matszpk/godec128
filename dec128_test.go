/*
 * dec128_test.go - main fixed decimal int128 routines
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
 
 package godec128
 
 import (
    "fmt"
    "strconv"
    "testing"
)

func getPanicInt2(f func(), paniced *bool, panicStr *string) {
    defer func() {
        if x:=recover(); x!=nil {
            *paniced = true
            *panicStr = fmt.Sprint(x)
        }
    }()
    f() // call
}

func getPanic2(f func()) (bool, string) {
    paniced := false
    panicStr := ""
    getPanicInt2(f, &paniced, &panicStr)
    return paniced, panicStr
}
 
 type UDec128TC struct {
    a, b UDec128
    expected UDec128
}
 
 func TestUDec128Add(t *testing.T) {
     testCases := []UDec128TC {
        UDec128TC{ UDec128{ 2454, 3421 }, UDec128{ 78731, 831 },
                UDec128{ 81185, 4252 } },
        UDec128TC{ UDec128{ 0xffffffffffff1001, 0x2442 }, UDec128{ 0xf003, 0xa8bc },
                UDec128{ 0x4, 0xccff } },
    }
    for i, tc := range testCases {
        a, b := tc.a, tc.b
        result := tc.a.Add(tc.b)
        if tc.expected!=result {
            t.Errorf("Result mismatch: %d: %v+%v->%v!=%v",
                     i, tc.a, tc.b, tc.expected, result)
        }
        if tc.a!=a || tc.b!=b {
            t.Errorf("Argument has been modified: %d: %v,%v!=%v,%v",
                     i, a, b, tc.a, tc.b)
        }
    }
 }
 
 func TestUDec128Sub(t *testing.T) {
    testCases := []UDec128TC {
        UDec128TC{ UDec128{ 81185, 4252 }, UDec128{ 2454, 3421 },
                UDec128{ 78731, 831 } },
        UDec128TC{ UDec128{ 0x4, 0xccff }, UDec128{ 0xffffffffffff1001, 0x2442 },
                UDec128{ 0xf003, 0xa8bc } },
    }
    for i, tc := range testCases {
        a, b := tc.a, tc.b
        result := tc.a.Sub(tc.b)
        if tc.expected!=result {
            t.Errorf("Result mismatch: %d: %v-%v->%v!=%v",
                     i, tc.a, tc.b, tc.expected, result)
        }
        if tc.a!=a || tc.b!=b {
            t.Errorf("Argument has been modified: %d: %v,%v!=%v,%v",
                     i, a, b, tc.a, tc.b)
        }
    }
}

type UDec128CTC struct {
    a, b UDec128
    c uint64
    expected UDec128
    expC uint64
}

func TestUDec128AddC(t *testing.T) {
    testCases := []UDec128CTC {
        UDec128CTC { UDec128{ 8481, 7754 }, UDec128{ 1121, 5531 }, 0,
            UDec128{ 9602, 13285 }, 0 },
        UDec128CTC { UDec128{ 8481, 7754 }, UDec128{ 1121, 5531 }, 1,
            UDec128{ 9603, 13285 }, 0 },
        UDec128CTC { UDec128{ 0xfffffffffffffffe, 7754 }, UDec128{ 1, 5531 }, 1,
            UDec128{ 0, 13286 }, 0 },
        UDec128CTC { UDec128{ 0xfffffffffffffffd, 7754 }, UDec128{ 1, 5531 }, 1,
            UDec128{ 0xffffffffffffffff, 13285 }, 0 },
        UDec128CTC { UDec128{ 0xffffffffffffff22, 0xfffffffffffffffe },
            UDec128{ 0xde, 1 }, 0, UDec128{ 0, 0 }, 1 },
        UDec128CTC { UDec128{ 0xffffffffffffff25, 0xfffffffffffffffe },
            UDec128{ 0xde, 2 }, 0, UDec128{ 3, 1 }, 1 },
        UDec128CTC { UDec128{ 0xffffffffffffff25, 0xfffffffffffffffe },
            UDec128{ 0xd1, 3 }, 0, UDec128{ 0xfffffffffffffff6, 1 }, 1 },
    }
    for i, tc := range testCases {
        a, b := tc.a, tc.b
        result, resultC := tc.a.AddC(tc.b, tc.c)
        if tc.expected!=result || tc.expC!=resultC {
            t.Errorf("Result mismatch: %d: addc(%v,%v,%v)->%v,%v!=%v,%v",
                     i, tc.a, tc.b, tc.c, tc.expected, tc.expC, result, resultC)
        }
        if tc.a!=a || tc.b!=b {
            t.Errorf("Argument has been modified: %d: %v,%v!=%v,%v",
                     i, a, b, tc.a, tc.b)
        }
    }
}

func TestUDec128SubB(t *testing.T) {
    testCases := []UDec128CTC {
        UDec128CTC{ UDec128{ 81183, 4252 }, UDec128{ 2454, 3421 }, 0,
                UDec128{ 78729, 831 }, 0 },
        UDec128CTC{ UDec128{ 81185, 4252 }, UDec128{ 2454, 3421 }, 1,
                UDec128{ 78730, 831 }, 0 },
        UDec128CTC{ UDec128{ 0x4, 0xccff }, UDec128{ 0xffffffffffff1001, 0x2442 }, 1,
                UDec128{ 0xf002, 0xa8bc }, 0 },
        UDec128CTC{ UDec128{ 81185, 4252 }, UDec128{ 81183, 4253 }, 0,
                UDec128{ 2 , 0xffffffffffffffff }, 1 },
        UDec128CTC{ UDec128{ 81185, 4252 }, UDec128{ 81187, 4253 }, 0,
                UDec128{ 0xfffffffffffffffe, 0xfffffffffffffffe }, 1 },
    }
    for i, tc := range testCases {
        a, b := tc.a, tc.b
        result, resultC := tc.a.SubB(tc.b, tc.c)
        if tc.expected!=result || tc.expC!=resultC {
            t.Errorf("Result mismatch: %d: subb(%v,%v,%v)->%v,%v!=%v,%v",
                     i, tc.a, tc.b, tc.c, tc.expected, tc.expC, result, resultC)
        }
        if tc.a!=a || tc.b!=b {
            t.Errorf("Argument has been modified: %d: %v,%v!=%v,%v",
                     i, a, b, tc.a, tc.b)
        }
    }
}

type UDec128CmpTC struct {
    a, b UDec128
    expected int
}

func TestUDec128Cmp(t *testing.T) {
    testCases := []UDec128CmpTC {
        UDec128CmpTC{ UDec128{ 3421, 2454 }, UDec128{ 831, 78731 }, -1 },
        UDec128CmpTC{ UDec128{ 6743, 6841 }, UDec128{ 7731121, 1212 }, 1 },
        UDec128CmpTC{ UDec128{ 1821, 33411 }, UDec128{ 589759892, 33411 }, -1 },
        UDec128CmpTC{ UDec128{ 5788219381, 33411 }, UDec128{ 954891, 33411 }, 1 },
        UDec128CmpTC{ UDec128{ 1231, 33411 }, UDec128{ 1231, 33411 }, 0 },
    }
    for i, tc := range testCases {
        a, b := tc.a, tc.b
        result := tc.a.Cmp(tc.b)
        if tc.expected!=result {
            t.Errorf("Result mismatch: %d: cmp(%v,%v)->%v!=%v",
                     i, tc.a, tc.b, tc.expected, result)
        }
        if tc.a!=a || tc.b!=b {
            t.Errorf("Argument has been modified: %d: %v,%v!=%v,%v",
                     i, a, b, tc.a, tc.b)
        }
    }
}

type UDec128MulTC struct {
    a, b UDec128
    precision uint
    rounding bool
    expected UDec128
}

func TestUDec128Mul(t *testing.T) {
    testCases := []UDec128MulTC {
        UDec128MulTC{ UDec128{ 0x840875a4212a9e42, 0x11310 },
                UDec128{ 0x3df9379d88970c7e, 0xc7 }, 8, false,
                UDec128{ 0xd3d0a477d6fda958, 0x23eaa838e89ce65c } },
        UDec128MulTC{ UDec128{ 0x840875a4212a9e43, 0x11310 }, // no rounding
                UDec128{ 0x3df9379d88970c7e, 0xc7 }, 8, false,
                UDec128{ 0xd3d0c5e538df353a, 0x23eaa838e89ce65c } },
        UDec128MulTC{ UDec128{ 0x840875a4212a9e43, 0x11310 }, // rounding
                UDec128{ 0x3df9379d88970c7e, 0xc7 }, 8, true,
                UDec128{ 0xd3d0c5e538df353b, 0x23eaa838e89ce65c } },
        UDec128MulTC{ UDec128{ 0x5d81bfe68a0b0c43, 0x65 },
                UDec128{ 0x089f625783250275, 0xb3cb }, 13, false,
                UDec128{ 0xc77d957642aa0de9, 0x7d3d5cc9dda } },
        UDec128MulTC{ UDec128{ 0x5d81bfe68a0b0c43, 0x65 },  // rounding
                UDec128{ 0x089f625783250275, 0xb3cb }, 13, true,
                UDec128{ 0xc77d957642aa0dea, 0x7d3d5cc9dda } },
    }
    for i, tc := range testCases {
        a, b := tc.a, tc.b
        result := tc.a.Mul(tc.b, tc.precision, tc.rounding)
        if tc.expected!=result {
            t.Errorf("Result mismatch: %d: mul(%v,%v,%v,%v)->%v!=%v",
                     i, tc.a, tc.b, tc.precision, tc.rounding, tc.expected, result)
        }
        if tc.a!=a || tc.b!=b {
            t.Errorf("Argument has been modified: %d: %v,%v!=%v,%v",
                     i, a, b, tc.a, tc.b)
        }
    }
}

type UDec128_64TC struct {
    a UDec128
    b uint64
    expected UDec128
}

func TestUDec128Mul64(t *testing.T) {
    testCases := []UDec128_64TC {
        UDec128_64TC { UDec128{ 0xc9baa109a40baa11, 0x384b9a928941ac3 },
            0x1839b9af9dc021, UDec128{ 0x6ac740f8d07aac31, 0x4e2dc743fec47ca9 } },
        UDec128_64TC { UDec128{ 0x2327f0eac961980e, 0x49f0f9 },
            0x11f82bb55bf, UDec128{ 0x77f8d53cd5871872, 0x530ae9c8b7cb9049 } },
    }
    for i, tc := range testCases {
        a, b := tc.a, tc.b
        result := tc.a.Mul64(tc.b)
        if tc.expected!=result {
            t.Errorf("Result mismatch: %d: %v*%v->%v!=%v",
                     i, tc.a, tc.b, tc.expected, result)
        }
        if tc.a!=a || tc.b!=b {
            t.Errorf("Argument has been modified: %d: %v,%v!=%v,%v",
                     i, a, b, tc.a, tc.b)
        }
    }
}

type UDec128MulFTC struct {
    a, b UDec128
    expectedLo, expectedHi UDec128
}

func TestUDec128MulFull(t *testing.T) {
    testCases := []UDec128MulFTC {
        UDec128MulFTC { UDec128{ 0xa0a59e0cd5640249, 0x5ff18c5e354dd456 },
            UDec128{ 0x4ddec0edfcc8c414, 0xadf9e6b9046f6ea3 },
            UDec128{ 0xf5a23257e29811b4, 0x89c07fdabef4588c },
            UDec128{ 0xd8d0c5c68299cf33, 0x4133e4458cfc0e8e } },
        UDec128MulFTC { UDec128{ 0xffffffffffffffff, 0xffffffffffffffff },
            UDec128{ 0xfffffffffffffffd, 0xffffffffffffffff },
            UDec128{ 3, 0 },
            UDec128{ 0xfffffffffffffffc, 0xffffffffffffffff } },
    }
    for i, tc := range testCases {
        a, b := tc.a, tc.b
        result, resultLo := tc.a.MulFull(tc.b)
        if tc.expectedHi!=result || tc.expectedLo!=resultLo {
            t.Errorf("Result mismatch: %d: mulfull(%v,%v)->%v,%v!=%v,%v",
                     i, tc.a, tc.b, tc.expectedLo, tc.expectedHi, resultLo, result)
        }
        if tc.a!=a || tc.b!=b {
            t.Errorf("Argument has been modified: %d: %v,%v!=%v,%v",
                     i, a, b, tc.a, tc.b)
        }
    }
}

type UDec128ShTC struct {
    a UDec128
    b uint
    expected UDec128
}

func TestUDec128Shl(t *testing.T) {
    testCases := []UDec128ShTC {
        UDec128ShTC { UDec128{ 0x62b71430f1765e8f, 0xb5ed145b3920ca5a }, 3,
            UDec128{ 0x15b8a1878bb2f478, 0xaf68a2d9c90652d3 } },
        UDec128ShTC { UDec128{ 0x62b71430f1765e8f, 0xb5ed145b3920ca5a }, 11,
            UDec128{ 0xb8a1878bb2f47800, 0x68a2d9c90652d315 } },
        UDec128ShTC { UDec128{ 0x62b7ac5532325e8f, 0xc5ed145b3920ca5a }, 0,
            UDec128{ 0x62b7ac5532325e8f, 0xc5ed145b3920ca5a } },
        UDec128ShTC { UDec128{ 0xf621e52aaa8b880c, 0xb4283ce0fd8464e2 }, 73,
            UDec128{ 0, 0x43ca555517101800 } },
        UDec128ShTC { UDec128{ 0xf621e52aaa8b880c, 0xb4283ce0fd8464e2 }, 64,
            UDec128{ 0, 0xf621e52aaa8b880c } },
    }
    for i, tc := range testCases {
        a, b := tc.a, tc.b
        result := tc.a.Shl(tc.b)
        if tc.expected!=result {
            t.Errorf("Result mismatch: %d:%v<<%v->%v!=%v",
                     i, tc.a, tc.b, tc.expected, result)
        }
        if tc.a!=a || tc.b!=b {
            t.Errorf("Argument has been modified: %d: %v,%v!=%v,%v",
                     i, a, b, tc.a, tc.b)
        }
    }
}

func TestUDec128Shr(t *testing.T) {
    testCases := []UDec128ShTC {
        UDec128ShTC { UDec128{ 0xeebbd1b847efeefa, 0x1f27b7128996878e }, 11,
            UDec128{ 0xf1ddd77a3708fdfd, 0x3e4f6e25132d0 } },
        UDec128ShTC { UDec128{ 0xecabd1b847efe63a, 0x1f27b7523196878f }, 0,
            UDec128{ 0xecabd1b847efe63a, 0x1f27b7523196878f } },
        UDec128ShTC { UDec128{ 0xf4f393b4762c797a, 0x51c18de532f49530 }, 82,
            UDec128{ 0x147063794cbd, 0 } },
        UDec128ShTC { UDec128{ 0xadd45555288f694c, 0x2b2e0d6f95ff2df1 }, 64,
            UDec128{ 0x2b2e0d6f95ff2df1, 0 } },
    }
    for i, tc := range testCases {
        a, b := tc.a, tc.b
        result := tc.a.Shr(tc.b)
        if tc.expected!=result {
            t.Errorf("Result mismatch: %d: %v>>%v->%v!=%v",
                     i, tc.a, tc.b, tc.expected, result)
        }
        if tc.a!=a || tc.b!=b {
            t.Errorf("Argument has been modified: %d: %v,%v!=%v,%v",
                     i, a, b, tc.a, tc.b)
        }
    }
}

type UDec128Div64TC struct {
    a UDec128
    b uint64
    expected UDec128
}

func TestUDec128Div64(t *testing.T) {
    testCases := []UDec128Div64TC {
        UDec128Div64TC { UDec128{ 0x0d362b7e0421d339, 0xbb09d477baa0 },
            0x6afcb5c6af1e507b, UDec128{ 492083670228144, 0 } },
        UDec128Div64TC { UDec128{ 0x0bc4f2ea7ec06c3f, 0x7bdcd02be78fe },
            0x3e2dc3dd417, UDec128{ 0xf6491fcb9513612d, 0x1fd } },
        UDec128Div64TC { UDec128{ 0, 1<<55 }, 1<<55, UDec128{ 0, 1 } },
        UDec128Div64TC { UDec128{ 0, 1<<62 }, 1<<55, UDec128{ 0, 128 } },
        UDec128Div64TC { UDec128{ 0x2e9700d1e595b258, 0x34a67968e5a },
            0xc23b96121, UDec128{ 0x64b6c9b6ee122e0c, 0x45 } },
        UDec128Div64TC { UDec128{ 55, 0 }, 7, UDec128{ 7, 0 } },
        // no remainder
        UDec128Div64TC { UDec128{ 0x0f2b92f72757046a, 0x15b807b7564a },
            0x26b380a13ca, UDec128{ 0xfaa679c50cd8d211, 0x8 } },
    }
    for i, tc := range testCases {
        a, b := tc.a, tc.b
        result:= tc.a.Div64(tc.b)
        if tc.expected!=result {
            t.Errorf("Result mismatch: %d: %v/%v->%v!=%v",
                     i, tc.a, tc.b, tc.expected, result)
        }
        if tc.a!=a || tc.b!=b {
            t.Errorf("Argument has been modified: %d: %v,%v!=%v,%v",
                     i, a, b, tc.a, tc.b)
        }
    }
}

type UDec128DivFTC struct {
    alo, ahi UDec128
    b UDec128
    expected UDec128
}

func TestUDec128DivFull(t *testing.T) {
    testCases := []UDec128DivFTC {
        UDec128DivFTC{ UDec128{ 0xa168b431ea4cbf25, 0xeeaf8afeafe15bf3 }, // alo
            UDec128{ 0x79da7cfc64734fc8, 0x1ae093566b591f }, // ahi
            UDec128{ 0x64611073ad67885c, 0x159b7addc721d10f }, // b
            UDec128{ 0x71e0ef1e6710ea31, 0x13e6fd8cef95977 } }, // quo
        UDec128DivFTC{ UDec128{ 0xc9a7d6e2cc4a9fe1, 0x7c5f7c4fe1dd3975 }, // alo
            UDec128{ 0x78c86ab5339b57fc, 0xaa9ea603a6ff1 }, // ahi
            UDec128{ 0x8d4959f4e6d39704, 0x17b4ad5d2b7537 }, // b
            UDec128{ 0x106a3b20e0f77e82, 0x73288478235baedf } }, // quo
        UDec128DivFTC{ UDec128{ 0xad1b0bef418b04f3, 0xad386b96ec18a75d }, // alo
            UDec128{ 0x3c179a833f04, 0 }, // ahi
            UDec128{ 0x448ab60d06e16d71, 0x21277fb3c975915 }, // b
            UDec128{ 0x1d00017916c509, 0 } }, // quo
        UDec128DivFTC{ UDec128{ 0xfe846594f784bcc1, 0xf3abd28b98484862 }, // alo
            UDec128{ 0xd3e91d7d4a, 0 }, // ahi
            UDec128{ 0x1725a5b765d6df45, 0x251135 }, // b
            UDec128{ 0x978eaa37efa35277, 0x5b788 } }, // quo
    }
    for i, tc := range testCases {
        alo, ahi, b := tc.alo, tc.ahi, tc.b
        result := UDec128DivFull(tc.ahi, tc.alo, tc.b)
        if tc.expected!=result {
            t.Errorf("Result mismatch: %d: (%v,%v)/%v->%v!=%v",
                     i, tc.alo, tc.ahi, tc.b, tc.expected, result)
        }
        if tc.alo!=alo || tc.ahi!=ahi || tc.b!=b {
            t.Errorf("Argument has been modified: %d: %v,%v,%v!=%v,%v,%v",
                     i, alo, ahi, b, tc.alo, tc.ahi, tc.b)
        }
    }
    
    paniced, panicStr := getPanic2(func() {
        UDec128DivFull(UDec128{ 0x54cd83b46f259de9, 0x213a9ec7 }, UDec128{},
                       UDec128{ 0x54cd83b46f259de9, 0x213a9ec7 })
    })
    if !paniced || panicStr!="Divide overflow" {
        t.Errorf("Unexpected panic: %v,%v", paniced, panicStr)
    }
    paniced, panicStr = getPanic2(func() {
        UDec128DivFull(UDec128{ 0x54cd834632566de9, 0x213a9ec7545 }, UDec128{},
                       UDec128{ 0x54cd83111f259663, 0x213a9ec7 })
    })
    if !paniced || panicStr!="Divide overflow" {
        t.Errorf("Unexpected panic: %v,%v", paniced, panicStr)
    }
    paniced, panicStr = getPanic2(func() {
        UDec128DivFull(UDec128{ 0x54cd834632566de9, 0x213a9ec7545 }, UDec128{},
                       UDec128{})
    })
    if !paniced || panicStr!="Divide by zero" {
        t.Errorf("Unexpected panic: %v,%v", paniced, panicStr)
    }
}

type UDec128DivTC struct {
    a, b UDec128
    precision uint
    expected UDec128
}

func TestUDec128Div(t *testing.T) {
    testCases := []UDec128DivTC {
        UDec128DivTC { UDec128{ 0x29d774b64027d71c, 0x50339e89 },
            UDec128{ 0xe1320b466aa1ee71, 0x9c }, 13, UDec128{ 0xa64cfe4e65832020, 0x4 } },
        UDec128DivTC { UDec128{ 0xaea112fccc354d11, 0x46b7da4 },
            UDec128{ 0xc2fea748532c9056, 0x4b30de }, 10, UDec128{ 0x2309736671, 0 } },
    }
    for i, tc := range testCases {
        a, b := tc.a, tc.b
        result := tc.a.Div(tc.b, tc.precision)
        if tc.expected!=result {
            t.Errorf("Result mismatch: %d: div(%v,%v,%v)->%v!=%v",
                     i, tc.a, tc.b, tc.precision, tc.expected, result)
        }
        if tc.a!=a || tc.b!=b {
            t.Errorf("Argument has been modified: %d: %v,%v!=%v,%v",
                     i, a, b, tc.a, tc.b)
        }
    }
}

type UDec128FmtTC struct {
    a UDec128
    precision uint
    trimZeroes bool
    expected string
}

type UDec128Fmt2TC struct {
    a UDec128
    precision uint
    dispPrecision uint
    trimZeroes bool
    expected string
}

func TestUDec128Format(t *testing.T) {
    testCases := []UDec128FmtTC {
        UDec128FmtTC{ UDec128{ 0x5f75348b0131b3af, 0xb3af0f }, 15, false,
            "217224419425.143693331510191" },
        UDec128FmtTC{ UDec128{ 0x5f75348b0131b3af, 0xb3af0f }, 10, false,
            "21722441942514369.3331510191" },
        UDec128FmtTC{ UDec128{ 0x5f75348b0131b3b8, 0xb3af0f }, 15, false,
            "217224419425.143693331510200" },
        UDec128FmtTC{ UDec128{ 0x5f75348b0131b3b8, 0xb3af0f }, 15, true,
            "217224419425.1436933315102" },
        UDec128FmtTC{ UDec128{ 0x5f75348b0131b2f0, 0xb3af0f }, 15, false,
            "217224419425.143693331510000" },
        UDec128FmtTC{ UDec128{ 0x5f75348b0131b2f0, 0xb3af0f }, 15, true,
            "217224419425.14369333151" },
        UDec128FmtTC{ UDec128{ 1984593924556, 0 }, 15, false,
            "0.001984593924556" },
        UDec128FmtTC{ UDec128{ 1984593924560, 0 }, 15, false,
            "0.001984593924560" },
        UDec128FmtTC{ UDec128{ 1984593924560, 0 }, 15, true,
            "0.00198459392456" },
        UDec128FmtTC{ UDec128{ 1984593924000, 0 }, 15, false,
            "0.001984593924000" },
        UDec128FmtTC{ UDec128{ 1984593924000, 0 }, 15, true,
            "0.001984593924" },
        UDec128FmtTC{ UDec128{ 0, 0 }, 15, true, "0.0" },
        UDec128FmtTC{ UDec128{ 1, 0 }, 15, false, "0.000000000000001" },
        UDec128FmtTC{ UDec128{ 3211984593924556, 0 }, 15, false,
            "3.211984593924556" },
        UDec128FmtTC{ UDec128{ 33000000000000000, 0 }, 15, false,
            "33.000000000000000" },
        UDec128FmtTC{ UDec128{ 33000000000000000, 0 }, 15, true,
            "33." },
        UDec128FmtTC{ UDec128{ 33400000000000000, 0 }, 15, true,
            "33.4" },
        UDec128FmtTC{ UDec128{ 33000400000000000, 0 }, 15, true,
            "33.0004" },
        // zero digits after comma
        UDec128FmtTC{ UDec128{ 0x5f75348b0131b3af, 0xb3af0f }, 0, false,
            "217224419425143693331510191" },
    }
    for i, tc := range testCases {
        a := tc.a
        result := tc.a.Format(tc.precision, tc.trimZeroes)
        if tc.expected!=result {
            t.Errorf("Result mismatch: %d: fmt(%v)->%v!=%v",
                     i, tc.a, tc.expected, result)
        }
        if tc.a!=a {
            t.Errorf("Argument has been modified: %d: %v!=%v", i, a, tc.a)
        }
        resultBytes := tc.a.FormatBytes(tc.precision, tc.trimZeroes)
        if tc.expected!=string(resultBytes) {
            t.Errorf("Result mismatch: %d: fmtBytes(%v)->%v!=%v",
                     i, tc.a, tc.expected, string(resultBytes))
        }
        if tc.a!=a {
            t.Errorf("Argument has been modified: %d: %v!=%v", i, a, tc.a)
        }
    }
    
    testCases2 := []UDec128Fmt2TC {
        UDec128Fmt2TC{ UDec128{ 0x5f75348b0131b3af, 0xb3af0f }, 15, 17, false,
            "217224419425.14369333151019100" },
        UDec128Fmt2TC{ UDec128{ 0x5f75348b0131b2f0, 0xb3af0f }, 15, 17, false,
            "217224419425.14369333151000000" },
        UDec128Fmt2TC{ UDec128{ 0x5f75348b0131b3af, 0xb3af0f }, 15, 12, false,
            "217224419425.143693331510" },
        UDec128Fmt2TC{ UDec128{ 0x5f75348b0131b3af, 0xb3af0f }, 15, 10, false,
            "217224419425.1436933315" },
        UDec128Fmt2TC{ UDec128{ 0x5f75348b0131b3af, 0xb3af0f }, 15, 17, true,
            "217224419425.143693331510191" },
        UDec128Fmt2TC{ UDec128{ 0x5f75348b0131b3af, 0xb3af0f }, 15, 10, true,
            "217224419425.1436933315" },
        UDec128Fmt2TC{ UDec128{ 0x5f75348b0131b3af, 0xb3af0f }, 15, 12, true,
            "217224419425.14369333151" },
        UDec128Fmt2TC{ UDec128{ 0x5f75348b0131b2f0, 0xb3af0f }, 15, 12, true,
            "217224419425.14369333151" },
        UDec128Fmt2TC{ UDec128{ 0x5f75348b0131b2f0, 0xb3af0f }, 15, 13, true,
            "217224419425.14369333151" },
    }
    for i, tc := range testCases2 {
        a := tc.a
        result := tc.a.FormatNew(tc.precision, tc.dispPrecision, tc.trimZeroes)
        if tc.expected!=result {
            t.Errorf("Result mismatch: %d: fmt(%v)->%v!=%v",
                     i, tc.a, tc.expected, result)
        }
        if tc.a!=a {
            t.Errorf("Argument has been modified: %d: %v!=%v", i, a, tc.a)
        }
        resultBytes := tc.a.FormatNewBytes(tc.precision, tc.dispPrecision, tc.trimZeroes)
        if tc.expected!=string(resultBytes) {
            t.Errorf("Result mismatch: %d: fmtBytes(%v)->%v!=%v",
                     i, tc.a, tc.expected, string(resultBytes))
        }
        if tc.a!=a {
            t.Errorf("Argument has been modified: %d: %v!=%v", i, a, tc.a)
        }
    }
}

type UDec128ParseTC struct {
    str string
    precision uint
    rounding bool
    expected UDec128
    expError error
}

func TestUDec128Parse(t *testing.T) {
    testCases := []UDec128ParseTC {
        UDec128ParseTC{ "217224419425.143693331510191", 15, false,
            UDec128{ 0x5f75348b0131b3af, 0xb3af0f }, nil },
        UDec128ParseTC{ "217224419425.1436933315101915", 15, false,
            UDec128{ 0x5f75348b0131b3af, 0xb3af0f }, nil },
        UDec128ParseTC{ "217224419425.143693331510191999", 15, false,
            UDec128{ 0x5f75348b0131b3af, 0xb3af0f }, nil },
        UDec128ParseTC{ "217224419425.1436933315101915", 15, true,
            UDec128{ 0x5f75348b0131b3b0, 0xb3af0f }, nil },
        UDec128ParseTC{ "217224419425.1436933315101", 15, false,
            UDec128{ 0x5f75348b0131b354, 0xb3af0f }, nil },
        UDec128ParseTC{ "39428394592112", 10, false,
            UDec128{ 0x2cf00161e4efc000, 0x537e }, nil },
        UDec128ParseTC{ "348943892891898938943893434921", 11, false,
            UDec128{}, strconv.ErrRange },
        UDec128ParseTC{ "0.001984593924556", 15, false,
            UDec128{ 1984593924556, 0 }, nil },
        UDec128ParseTC{ ".0019845939245565", 15, false,
            UDec128{ 1984593924556, 0 }, nil },
        UDec128ParseTC{ ".0019845939245565", 15, true,
            UDec128{ 1984593924557, 0 }, nil },
        UDec128ParseTC{ "0.001984593924560", 15, false,
            UDec128{ 1984593924560, 0 }, nil },
        UDec128ParseTC{ ".001984593924560", 15, false,
            UDec128{ 1984593924560, 0 }, nil },
        UDec128ParseTC{ "0.00198459392456", 15, false,
            UDec128{ 1984593924560, 0 }, nil },
        UDec128ParseTC{ ".00198459392456", 15, false,
            UDec128{ 1984593924560, 0 }, nil },
        UDec128ParseTC{ ".001984593924", 15, false,
            UDec128{ 1984593924000, 0 }, nil },
        UDec128ParseTC{ "0.201984593924556", 15, false,
            UDec128{ 201984593924556, 0 }, nil },
        UDec128ParseTC{ ".30198459392456", 15, false,
            UDec128{ 301984593924560, 0 }, nil },
        UDec128ParseTC{ "0.0", 10, false, UDec128{}, nil },
        UDec128ParseTC{ "0", 10, false, UDec128{}, nil },
        UDec128ParseTC{ "0.", 10, false, UDec128{}, nil },
        UDec128ParseTC{ ".0", 10, false, UDec128{}, nil },
        UDec128ParseTC{ "217224419425.143693331510191e0", 15, false,
            UDec128{ 0x5f75348b0131b3af, 0xb3af0f }, nil },
        UDec128ParseTC{ "21722441942.5143693331510191e1", 15, false,
            UDec128{ 0x5f75348b0131b3af, 0xb3af0f }, nil },
        UDec128ParseTC{ "21722441942.5143693331510191ee1", 15, false,
            UDec128{}, strconv.ErrSyntax },
        UDec128ParseTC{ "2172244194.25143693331510191e2", 15, false,
            UDec128{ 0x5f75348b0131b3af, 0xb3af0f }, nil },
        UDec128ParseTC{ "2172244.19425143693331510192e5", 15, false,
            UDec128{ 0x5f75348b0131b3b0, 0xb3af0f }, nil },
        UDec128ParseTC{ "2172244194251.43693331510191e-1", 15, false,
            UDec128{ 0x5f75348b0131b3af, 0xb3af0f }, nil },
        UDec128ParseTC{ "217224419425143693.331510190e-6", 15, false,
            UDec128{ 0x5f75348b0131b3ae, 0xb3af0f }, nil },
        UDec128ParseTC{ "0.01984593924556e-1", 15, false,
            UDec128{ 1984593924556, 0 }, nil },
        UDec128ParseTC{ ".01984593924556e-1", 15, false,
            UDec128{ 1984593924556, 0 }, nil },
        UDec128ParseTC{ "0.1984593924556e-2", 15, false,
            UDec128{ 1984593924556, 0 }, nil },
        UDec128ParseTC{ "00.1984593924556e-2", 15, false,
            UDec128{ 1984593924556, 0 }, nil },
        UDec128ParseTC{ ".1984593924556e-2", 15, false,
            UDec128{ 1984593924556, 0 }, nil },
        UDec128ParseTC{ "1.984593924556e-3", 15, false,
            UDec128{ 1984593924556, 0 }, nil },
        UDec128ParseTC{ "12e3", 15, false, UDec128{ 12000000000000000000, 0 }, nil },
        UDec128ParseTC{ "12.e3", 15, false, UDec128{ 12000000000000000000, 0 }, nil },
        UDec128ParseTC{ "12.77e3", 15, false, UDec128{ 12770000000000000000, 0 }, nil },
        UDec128ParseTC{ "0.0e0", 10, false, UDec128{}, nil },
        UDec128ParseTC{ "0.0e1", 10, false, UDec128{}, nil },
        UDec128ParseTC{ ".0e1", 10, false, UDec128{}, nil },
        UDec128ParseTC{ "0.e1", 10, false, UDec128{}, nil },
        UDec128ParseTC{ "0.0e3", 10, false, UDec128{}, nil },
        UDec128ParseTC{ ".0e3", 10, false, UDec128{}, nil },
        UDec128ParseTC{ "0.e3", 10, false, UDec128{}, nil },
        UDec128ParseTC{ "0.0e-1", 10, false, UDec128{}, nil },
        UDec128ParseTC{ ".0e-1", 10, false, UDec128{}, nil },
        UDec128ParseTC{ "0.e-1", 10, false, UDec128{}, nil },
        UDec128ParseTC{ "0.0e-3", 10, false, UDec128{}, nil },
        UDec128ParseTC{ ".0e-3", 10, false, UDec128{}, nil },
        UDec128ParseTC{ "0.e-3", 10, false, UDec128{}, nil },
        UDec128ParseTC{ "12344", 0, false, UDec128{ 12344, 0 }, nil },
        UDec128ParseTC{ "12344.", 0, false, UDec128{ 12344, 0 }, nil },
        UDec128ParseTC{ "12344.0000", 0, false, UDec128{ 12344, 0 }, nil },
        UDec128ParseTC{ "12344.7000", 0, false, UDec128{ 12344, 0 }, nil },
        UDec128ParseTC{ "12344.7000", 0, true, UDec128{ 12345, 0 }, nil },
    }
    for i, tc := range testCases {
        result, err := ParseUDec128(tc.str, tc.precision, tc.rounding)
        if tc.expected!=result || tc.expError!=err {
            t.Errorf("Result mismatch: %d: parse(%v)->%v,%v!=%v,%v",
                     i, tc.str, tc.expected, tc.expError, result, err)
        }
        result, err = ParseUDec128Bytes([]byte(tc.str), tc.precision, tc.rounding)
        if tc.expected!=result || tc.expError!=err {
            t.Errorf("Result mismatch: %d: parse(%v)->%v,%v!=%v,%v",
                     i, tc.str, tc.expected, tc.expError, result, err)
        }
    }
}

type UDec128ToFloat64TC struct {
    value UDec128
    precision uint
    expected float64
}

func TestUDec128ToFloat64(t *testing.T) {
    testCases := []UDec128ToFloat64TC{
        UDec128ToFloat64TC{ UDec128{ 0, 0 }, 11, 0.0 },
        UDec128ToFloat64TC{ UDec128{ 1, 0 }, 11, 1.0*1e-11 },
        UDec128ToFloat64TC{ UDec128{ 54930201, 0 }, 11, 54930201.0*1e-11 },
        UDec128ToFloat64TC{ UDec128{ 85959028918918968, 0 }, 0,
                    85959028918918968.0 },
        UDec128ToFloat64TC{ UDec128{ 85959028918918968, 0 }, 11,
                    85959028918918968.0*1e-11 },
        UDec128ToFloat64TC{ UDec128{ 85959028918918968, 0 }, 17,
                    0.8595902891891898 },
        UDec128ToFloat64TC{ UDec128{ 16346246572275455745, 10277688839402 }, 11,
                    189589895689685989335661129029377.0*1e-11 },
        UDec128ToFloat64TC{ UDec128{ 0xffffffffffffffff, 0xffffffffffffffff }, 11,
                    340282366920938463463374607431768211455.0*1e-11 },
    }
    for i, tc := range testCases {
        result := tc.value.ToFloat64(tc.precision)
        if tc.expected!=result {
            t.Errorf("Result mismatch: %d: tofloat64(%v,%v)->%v!=%v",
                     i, tc.value, tc.precision, tc.expected, result)
        }
    }
}

type Float64ToUDec128TC struct {
    value float64
    precision uint
    expected UDec128
    expError error
}

func TestFloat64ToUDec128(t *testing.T) {
    testCases := []Float64ToUDec128TC{
        Float64ToUDec128TC{ 0.0, 0, UDec128{ 0, 0 }, nil },
        Float64ToUDec128TC{ 1.0, 0, UDec128{ 1, 0 }, nil },
        Float64ToUDec128TC{ 1.7, 0, UDec128{ 1, 0 }, nil },
        Float64ToUDec128TC{ 145645677.18, 0, UDec128{ 145645677, 0 }, nil },
        Float64ToUDec128TC{ 3145645677.778, 0, UDec128{ 3145645677, 0 }, nil },
        Float64ToUDec128TC{ 187923786919586921.0, 0,
            UDec128{ 187923786919586912, 0 }, nil },
        Float64ToUDec128TC{ 11792378691958692154.0, 0,
            UDec128{ 11792378691958691840, 0 }, nil },
        Float64ToUDec128TC{ 26858969188828978177.0, 0,
            UDec128{ 8412225115119427584, 1 }, nil },
        Float64ToUDec128TC{ 145645677.18, 3, UDec128{ 145645677180, 0 }, nil },
        Float64ToUDec128TC{ 58590303.45539292211, 11,
            UDec128{ 0x514f750e8a1a8c00, 0 }, nil },
    }
    for i, tc := range testCases {
        result, err := Float64ToUDec128(tc.value, tc.precision)
        if tc.expected!=result || tc.expError!=err {
            t.Errorf("Result mismatch: %d: toudec128(%v)->%v,%v!=%v,%v",
                     i, tc.value, tc.expected, tc.expError, result, err)
        }
    }
}
