package vrtime

import (
    "math"
)

type Key1_type int64
type Key2_type int64

type Time struct {
    Key1 Key1_type
    Key2 Key2_type
}

func cmpKey1 (lhs Key1_type, rhs Key1_type) int {
    if lhs < rhs {
        return -1
    } else if lhs > rhs {
        return 1
    }
    return 0  
}

func cmpKey2 (lhs Key2_type, rhs Key2_type) int {
    if lhs < rhs {
        return -1
    } else if lhs > rhs {
        return 1
    }
    return 0  
}

func ZeroTime() Time {
    return Time{Key1:0,Key2:0}
}

func InfinityTime() Time {
    return Time{Key1: math.MaxInt64,Key2: math.MaxInt64}
}

func (s Time) LT (t Time) bool {
    ck1 := cmpKey1(s.Key1, t.Key1) 
    if ck1 == -1 {
        return true
    } else if ck1 == 1 {
        return false
    }
    return cmpKey2(s.Key2, t.Key2) == -1
}

func (s Time) GT (t Time) bool {
    ck1 := cmpKey1(s.Key1, t.Key1) 
    if ck1 == 1 {
        return true
    } else if ck1 == -1 {
        return false
    }
    return cmpKey2(s.Key2, t.Key2) == 1
}

func (s Time) EQ (t Time) bool {
    ck1 := cmpKey1(s.Key1, t.Key1) 
    if ck1 != 0 {
        return false 
    }
    return cmpKey2(s.Key2, t.Key2) == 0
}

func (s Time) LE (t Time) bool {
    return s.LE(t) || s.EQ(t)
}

func (s Time) GE (t Time) bool {
    return s.GE(t) || s.EQ(t)
}

func (s Time) NEQ (t Time) bool {
    return !s.EQ(t) 
}

func (s Time) Plus (t Time) Time {
    sum := Time{s.Key1+t.Key1, s.Key2+t.Key2}
    return sum
}
