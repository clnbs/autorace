package mathtool

import "testing"

func TestRandomIntBetween(t *testing.T) {
	for i := 0; i < 10; i++ {
		result := RandomIntBetween(5, 10)
		if result < 5 || result > 10 {
			t.Fatal("result should be between 5 and 10, got ", result)
		}
	}
	for i := 0; i < 10; i++ {
		result := RandomIntBetween(10, 5)
		if result < 5 || result > 10 {
			t.Fatal("result should be between 5 and 10, got ", result)
		}
	}
	for i := 0; i < 10; i++ {
		result := RandomIntBetween(-5, -10)
		if result < -10 || result > -5 {
			t.Fatal("result should be between -5 and -10, got ", result)
		}
	}
	for i := 0; i < 10; i++ {
		result := RandomIntBetween(-10, -5)
		if result < -10 || result > -5 {
			t.Fatal("result should be between -5 and -10, got ", result)
		}
	}
	for i := 0; i < 10; i++ {
		result := RandomIntBetween(-5, 10)
		if result < -5 || result > 10 {
			t.Fatal("result should be between -5 and 10, got ", result)
		}
	}
	for i := 0; i < 10; i++ {
		result := RandomIntBetween(10, -5)
		if result < -5 || result > 10 {
			t.Fatal("result should be between -5 and 10, got ", result)
		}
	}
	result := RandomIntBetween(5, 5)
	if result != 5 {
		t.Fatal("result should be equals to 5, got ", result)
	}
}

func TestRandomFloatBetween(t *testing.T) {
	for i := 0; i < 10; i++ {
		result := RandomFloatBetween(5.5, 10.5)
		if result < 5.5 || result > 10.5 {
			t.Fatal("result should be between 5 and 10, got ", result)
		}
	}
	for i := 0; i < 10; i++ {
		result := RandomFloatBetween(10.5, 5.5)
		if result < 5.5 || result > 10.5 {
			t.Fatal("result should be between 5 and 10, got ", result)
		}
	}
	for i := 0; i < 10; i++ {
		result := RandomFloatBetween(-5.5, -10.5)
		if result < -10.5 || result > -5.5 {
			t.Fatal("result should be between -5 and -10, got ", result)
		}
	}
	for i := 0; i < 10; i++ {
		result := RandomFloatBetween(-10.5, -5.5)
		if result < -10.5 || result > -5.5 {
			t.Fatal("result should be between -5 and -10, got ", result)
		}
	}
	for i := 0; i < 10; i++ {
		result := RandomFloatBetween(-5.5, 10.5)
		if result < -5.5 || result > 10.5 {
			t.Fatal("result should be between -5 and 10, got ", result)
		}
	}
	for i := 0; i < 10; i++ {
		result := RandomFloatBetween(10.5, -5.5)
		if result < -5.5 || result > 10.5 {
			t.Fatal("result should be between -5 and 10, got ", result)
		}
	}
	result := RandomFloatBetween(5.5, 5.5)
	if result != 5.5 {
		t.Fatal("result should be equals to 5.5, got ", result)
	}
}

func TestClampFloat64(t *testing.T) {
	x := 5.5
	result := ClampFloat64(x, 0.0, 5.0)
	if result != 5.0 {
		t.Fatal("result should be equal to 5.0 got", result)
	}
	result = ClampFloat64(x, 6.0, 10.0)
	if result != 6.0 {
		t.Fatal("result should be equal to 6.0 got", result)
	}
	result = ClampFloat64(x, 0.0, 10.0)
	if result != x {
		t.Fatal("result should be equal to", x, "got", result)
	}

	result = ClampFloat64(x, 5.0, 0.0)
	if result != 5.0 {
		t.Fatal("result should be equal to 5.0 got", result)
	}
	result = ClampFloat64(x, 10.0, 6.0)
	if result != 6.0 {
		t.Fatal("result should be equal to 6.0 got", result)
	}
	result = ClampFloat64(x, 10.0, 0.0)
	if result != x {
		t.Fatal("result should be equal to", x, "got", result)
	}
}

func TestClampInt(t *testing.T) {
	x := 5
	result := ClampInt(x, 0, 4)
	if result != 4 {
		t.Fatal("result should be equal to 4, got", result)
	}
	result = ClampInt(x, 6, 10)
	if result != 6 {
		t.Fatal("result should be equal to 6, got", result)
	}
	result = ClampInt(x, 0, 10)
	if result != x {
		t.Fatal("result shoub be equal to", x, "got", result)
	}

	result = ClampInt(x, 4, 0)
	if result != 4 {
		t.Fatal("result should be equal to 4, got", result)
	}
	result = ClampInt(x, 10, 6)
	if result != 6 {
		t.Fatal("result should be equal to 6, got", result)
	}
	result = ClampInt(x, 10, 0)
	if result != x {
		t.Fatal("result shoub be equal to", x, "got", result)
	}
}

func TestIsFloat64Between(t *testing.T) {
	x := 5.0
	if !IsFloat64Between(x, 0.0, 10.0) {
		t.Fatal(x, "should be between 0.0 and 10.0")
	}
	if !IsFloat64Between(x, 10.0, 0.0) {
		t.Fatal(x, "should be between 10.0 and 0.0")
	}
	if IsFloat64Between(x, 10.0, 20.0) {
		t.Fatal(x, "should not be between 10.0 and 20.0")
	}
	if IsFloat64Between(x, 20.0, 10.0) {
		t.Fatal(x, "should not be between 20.0 and 10.0")
	}
}

func TestIsIntBetween(t *testing.T) {
	x := 5
	if !IsIntBetween(x, 0, 10) {
		t.Fatal(x, "should be between 0 and 10")
	}
	if !IsIntBetween(x, 10, 0) {
		t.Fatal(x, "should be between 10 and 0")
	}
	if IsIntBetween(x, 10, 20) {
		t.Fatal(x, "should not be between 10 and 20")
	}
	if IsIntBetween(x, 20, 10) {
		t.Fatal(x, "should not be between 20 and 10")
	}
}
