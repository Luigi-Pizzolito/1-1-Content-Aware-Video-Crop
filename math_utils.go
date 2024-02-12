package main

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}

func pos(x int) int {
    if x < 0 {
        return 0
    }
    return x
}

func posF(x float64) float64 {
    if x < 0 {
        return 0
    }
    return x
}

func clamp(value, min, max int) int {
    if value < min {
        return min
    }
    if value > max {
        return max
    }
    return value
}

func clampF(value, min, max float64) float64 {
    if value < min {
        return min
    }
    if value > max {
        return max
    }
    return value
}