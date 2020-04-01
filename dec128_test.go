/*
 * dec128_test.go - main fixed decimal int128 routines
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
 
 package godec128
 
 import (
    "testing"
)
 
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

type UDec128_64TC struct {
    a UDec128
    b uint64
    expected UDec128
}

func TestUDec128Add64(t *testing.T) {
    testCases := []UDec128_64TC {
        UDec128_64TC{ UDec128{ 3454, 3421 }, 78731, UDec128{ 82185, 3421 } },
        UDec128_64TC{ UDec128{ 0xffffffffffff1001, 0x2446 }, 0xf003,
                UDec128{ 0x4, 0x2447 } },
    }
    for i, tc := range testCases {
        a, b := tc.a, tc.b
        result := tc.a.Add64(tc.b)
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

func TestUDec128Sub64(t *testing.T) {
    testCases := []UDec128_64TC {
        UDec128_64TC{ UDec128{ 81185, 9165 }, 2454, UDec128{ 78731, 9165 } },
        UDec128_64TC{ UDec128{ 0x5, 0xccff }, 0xffffffffffff2001,
                UDec128{ 0xe004, 0xccfe } },
    }
    for i, tc := range testCases {
        a, b := tc.a, tc.b
        result := tc.a.Sub64(tc.b)
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

type UDec128Mul struct {
    a, b UDec128
    tenPow uint
    
    expected UDec128
}
