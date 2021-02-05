/*
 * locale_test.go - tests for int128 routines
 *
 * godec128 - go dec128 (for 128-bit decimal fixed point) library
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
    "strconv"
    "testing"
)

type UDec128LocTC struct {
    lang string
    noSep1000 bool
    a UDec128
    precision uint
    trimZeroes bool
    expected string
}

func TestUDec128LocaleFormat(t *testing.T) {
    testCases := []UDec128LocTC {
        UDec128LocTC{ "af", false, UDec128{0xab54a98ceb1f0ad3, 0},
                10, false, "1 234 567 890,1234567891" },
        UDec128LocTC{ "af", false, UDec128{0xab54a98ceb1f0ad2, 0},
                10, false, "1 234 567 890,1234567890" },
        UDec128LocTC{ "af", false, UDec128{0xab54a98ceb1f0ad2, 0},
                10, true, "1 234 567 890,123456789" },
        UDec128LocTC{ "am", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1,234,567,890.1234567891" },
        UDec128LocTC{ "ar", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "١٬٢٣٤٬٥٦٧٬٨٩٠٫١٢٣٤٥٦٧٨٩١" },
        UDec128LocTC{ "az", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1.234.567.890,1234567891" },
        UDec128LocTC{ "bg", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1 234 567 890,1234567891" },
        UDec128LocTC{ "bn", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "১,২৩,৪৫,৬৭,৮৯০.১২৩৪৫৬৭৮৯১" },
        UDec128LocTC{ "ca", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1.234.567.890,1234567891" },
        UDec128LocTC{ "cs", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1 234 567 890,1234567891" },
        UDec128LocTC{ "da", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1.234.567.890,1234567891" },
        UDec128LocTC{ "de", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1.234.567.890,1234567891" },
        UDec128LocTC{ "el", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1.234.567.890,1234567891" },
        UDec128LocTC{ "en", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1,234,567,890.1234567891" },
        UDec128LocTC{ "es", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1.234.567.890,1234567891" },
        UDec128LocTC{ "et", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1 234 567 890,1234567891" },
        UDec128LocTC{ "fa", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "۱٬۲۳۴٬۵۶۷٬۸۹۰٫۱۲۳۴۵۶۷۸۹۱" },
        UDec128LocTC{ "fi", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1 234 567 890,1234567891" },
        UDec128LocTC{ "fil", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1,234,567,890.1234567891" },
        UDec128LocTC{ "fr", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1 234 567 890,1234567891" },
        UDec128LocTC{ "gu", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1,23,45,67,890.1234567891" },
        UDec128LocTC{ "he", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1,234,567,890.1234567891" },
        UDec128LocTC{ "hi", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1,23,45,67,890.1234567891" },
        UDec128LocTC{ "hr", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1.234.567.890,1234567891" },
        UDec128LocTC{ "hu", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1 234 567 890,1234567891" },
        UDec128LocTC{ "hy", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1 234 567 890,1234567891" },
        UDec128LocTC{ "id", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1.234.567.890,1234567891" },
        UDec128LocTC{ "is", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1.234.567.890,1234567891" },
        UDec128LocTC{ "it", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1.234.567.890,1234567891" },
        UDec128LocTC{ "ja", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1,234,567,890.1234567891" },
        UDec128LocTC{ "ka", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1 234 567 890,1234567891" },
        UDec128LocTC{ "kk", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1 234 567 890,1234567891" },
        UDec128LocTC{ "km", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1.234.567.890,1234567891" },
        UDec128LocTC{ "kn", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1,234,567,890.1234567891" },
        UDec128LocTC{ "ko", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1,234,567,890.1234567891" },
        UDec128LocTC{ "ky", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1 234 567 890,1234567891" },
        UDec128LocTC{ "lo", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1.234.567.890,1234567891" },
        UDec128LocTC{ "lt", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1 234 567 890,1234567891" },
        UDec128LocTC{ "lv", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1 234 567 890,1234567891" },
        UDec128LocTC{ "mk", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1.234.567.890,1234567891" },
        UDec128LocTC{ "ml", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1,23,45,67,890.1234567891" },
        UDec128LocTC{ "mn", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1,234,567,890.1234567891" },
        UDec128LocTC{ "mo", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1.234.567.890,1234567891" },
        UDec128LocTC{ "mr", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "१,२३,४५,६७,८९०.१२३४५६७८९१" },
        UDec128LocTC{ "ms", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1,234,567,890.1234567891" },
        UDec128LocTC{ "mul", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1,234,567,890.1234567891" },
        UDec128LocTC{ "my", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "၁,၂၃၄,၅၆၇,၈၉၀.၁၂၃၄၅၆၇၈၉၁" },
        UDec128LocTC{ "nb", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1 234 567 890,1234567891" },
        UDec128LocTC{ "ne", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "१,२३४,५६७,८९०.१२३४५६७८९१" },
        UDec128LocTC{ "nl", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1.234.567.890,1234567891" },
        UDec128LocTC{ "no", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1,234,567,890.1234567891" },
        UDec128LocTC{ "pa", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1,23,45,67,890.1234567891" },
        UDec128LocTC{ "pl", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1 234 567 890,1234567891" },
        UDec128LocTC{ "pt", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1.234.567.890,1234567891" },
        UDec128LocTC{ "ro", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1.234.567.890,1234567891" },
        UDec128LocTC{ "ru", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1 234 567 890,1234567891" },
        UDec128LocTC{ "sh", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1.234.567.890,1234567891" },
        UDec128LocTC{ "si", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1,234,567,890.1234567891" },
        UDec128LocTC{ "sk", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1 234 567 890,1234567891" },
        UDec128LocTC{ "sl", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1.234.567.890,1234567891" },
        UDec128LocTC{ "sq", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1 234 567 890,1234567891" },
        UDec128LocTC{ "sr", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1.234.567.890,1234567891" },
        UDec128LocTC{ "sv", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1 234 567 890,1234567891" },
        UDec128LocTC{ "sw", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1,234,567,890.1234567891" },
        UDec128LocTC{ "ta", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1,23,45,67,890.1234567891" },
        UDec128LocTC{ "te", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1,234,567,890.1234567891" },
        UDec128LocTC{ "th", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1,234,567,890.1234567891" },
        UDec128LocTC{ "tl", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1,234,567,890.1234567891" },
        UDec128LocTC{ "tn", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1,234,567,890.1234567891" },
        UDec128LocTC{ "tr", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1.234.567.890,1234567891" },
        UDec128LocTC{ "uk", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1 234 567 890,1234567891" },
        UDec128LocTC{ "ur", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1,234,567,890.1234567891" },
        UDec128LocTC{ "uz", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1 234 567 890,1234567891" },
        UDec128LocTC{ "vi", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1.234.567.890,1234567891" },
        UDec128LocTC{ "zh", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1,234,567,890.1234567891" },
        UDec128LocTC{ "zu", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1,234,567,890.1234567891" },
        UDec128LocTC{ "", false, UDec128{0xab54a98ceb1f0ad3,0}, 10,
                false, "1,234,567,890.1234567891" },
        UDec128LocTC{ "C", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1,234,567,890.1234567891" },
        UDec128LocTC{ "pl-PL", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1 234 567 890,1234567891" },
        UDec128LocTC{ "pl_PL", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1 234 567 890,1234567891" },
        UDec128LocTC{ "pl_PL.UTF-8", false, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1 234 567 890,1234567891" },
        // no separator 1000
        UDec128LocTC{ "af", true, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1234567890,1234567891" },
        UDec128LocTC{ "am", true, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1234567890.1234567891" },
        UDec128LocTC{ "ar", true, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "١٢٣٤٥٦٧٨٩٠٫١٢٣٤٥٦٧٨٩١" },
        UDec128LocTC{ "az", true, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1234567890,1234567891" },
        UDec128LocTC{ "bg", true, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1234567890,1234567891" },
        UDec128LocTC{ "bn", true, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "১২৩৪৫৬৭৮৯০.১২৩৪৫৬৭৮৯১" },
        UDec128LocTC{ "ca", true, UDec128{0xab54a98ceb1f0ad3,0},
                10, false, "1234567890,1234567891" },
    }
    for i, tc := range testCases {
        a := tc.a
        result := tc.a.LocaleFormat(tc.lang, tc.precision, tc.trimZeroes, tc.noSep1000)
        if tc.expected!=result {
            t.Errorf("Result mismatch: %d: fmt(%v,%s,%v,%v)->%v!=%v",
                     i, tc.a, tc.lang, tc.precision, tc.trimZeroes, tc.expected, result)
        }
        if tc.a!=a {
            t.Errorf("Argument has been modified: %d %s: %v!=%v", i, tc.lang, a, tc.a)
        }
        resultBytes := tc.a.LocaleFormatBytes(tc.lang, tc.precision,
                        tc.trimZeroes, tc.noSep1000)
        if tc.expected!=string(resultBytes) {
            t.Errorf("Result mismatch: %d: fmtBytes(%v,%s,%v,%v)->%v!=%v",
                     i, tc.a, tc.lang, tc.precision, tc.trimZeroes, tc.expected, result)
        }
        if tc.a!=a {
            t.Errorf("Argument has been modified: %d %s: %v!=%v", i, tc.lang, a, tc.a)
        }
    }
}

type UDec128LocParseTC struct {
    lang string
    str string
    precision uint
    rounding bool
    expected UDec128
    expError error
}

func TestUDec128LocaleParse(t *testing.T) {
    testCases := []UDec128LocParseTC {
        UDec128LocParseTC{ "en", "", 10, false, UDec128{}, strconv.ErrSyntax },
        UDec128LocParseTC{ "en", "1,234,567,890.1234567891", 10, false,
                UDec128{0xab54a98ceb1f0ad3,0}, nil },
        UDec128LocParseTC{ "en", "1,234,567,890.12345678915", 10, false,
                UDec128{0xab54a98ceb1f0ad3,0}, nil },
        UDec128LocParseTC{ "en", "1,234,567,890.12345678915", 10, true,
                UDec128{0xab54a98ceb1f0ad4,0}, nil },
        UDec128LocParseTC{ "en", "1,234,567890.1234567891", 10, false,
                UDec128{0xab54a98ceb1f0ad3,0}, nil },
        UDec128LocParseTC{ "pl", "1 234 567 890,1234567891", 10, false,
                UDec128{0xab54a98ceb1f0ad3,0}, nil },
        UDec128LocParseTC{ "pl", "1 234 567 890,1234567891", 10, false,
                UDec128{0xab54a98ceb1f0ad3,0}, nil },
        UDec128LocParseTC{ "bn", "১,২৩,৪৫,৬৭,৮৯০.১২৩৪৫৬৭৮৯১", 10, false,
                UDec128{0xab54a98ceb1f0ad3,0}, nil },
        UDec128LocParseTC{ "bn", "1,234,567890.1234567891", 10, false,
                UDec128{0xab54a98ceb1f0ad3,0}, nil },
    }
    for i, tc := range testCases {
        result, err := LocaleParseUDec128(tc.lang, tc.str, tc.precision, tc.rounding)
        if tc.expected!=result || tc.expError!=err {
            t.Errorf("Result mismatch: %d: parse(%v,%v,%v,%v)->%v,%v!=%v,%v",
                     i, tc.lang, tc.str, tc.precision, tc.rounding,
                     tc.expected, tc.expError, result, err)
        }
        result, err = LocaleParseUDec128Bytes(tc.lang, []byte(tc.str),
                                tc.precision, tc.rounding)
        if tc.expected!=result || tc.expError!=err {
            t.Errorf("Result mismatch: %d: parseBytes(%v,%v,%v,%v)->%v,%v!=%v,%v",
                     i, tc.lang, tc.str, tc.precision, tc.rounding,
                     tc.expected, tc.expError, result, err)
        }
    }
}

func BenchmarkUDec128LocaleFormatBytes(b *testing.B) {
    a := UDec128{ 7341542494928938945, 938491 }
    for i := 0; i < b.N; i++ {
        a.LocaleFormatBytes("pl", 8, false, false)
    }
}

func BenchmarkUDec128LocaleFormat(b *testing.B) {
    a := UDec128{ 7341542494928938945, 938491 }
    for i := 0; i < b.N; i++ {
        a.LocaleFormat("pl", 8, false, false)
    }
}
