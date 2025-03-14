package model

import (
	"sync"
	"testing"
)

func TestCopyRace(t *testing.T) {
	t.Parallel()
	a := map[string]int{"key1": 1, "key2": 2, "key3": 3, "key4": 4, "key5": 5, "key6": 6, "key7": 7, "key8": 8, "key9": 9, "key10": 10, "key11": 11, "key12": 12, "key13": 13, "key14": 14, "key15": 15,
		"key16": 16, "key17": 17, "key18": 18, "key19": 19, "key20": 20, "key21": 21, "key22": 22, "key23": 23, "key24": 24, "key25": 25, "key26": 26, "key27": 27}

	ast := &Assertion{
		Key:       "test_key",
		Value:     "test_value",
		Tokens:    []string{"token1", "token2"},
		Policy:    [][]string{{"policy1"}, {"policy2"}},
		PolicyMap: a,
	}

	var wg sync.WaitGroup
	const numRoutines = 100

	wg.Add(numRoutines)
	for i := 0; i < numRoutines; i++ {
		go func() {
			defer wg.Done()
			_ = ast.copy() // Просто вызываем copy(), не проверяя результат
		}()
	}
	wg.Wait()

	// Проверяем, что PolicyMap остался неизменным
	expectedPolicyMap := a
	if !compareMaps(expectedPolicyMap, ast.PolicyMap) {
		t.Errorf("PolicyMap was modified during concurrent access")
	}
}

func compareMaps(m1, m2 map[string]int) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v := range m1 {
		if v2, ok := m2[k]; !ok || v2 != v {
			return false
		}
	}
	return true
}

func BenchmarkCopy(b *testing.B) {
	ast := &Assertion{
		Key:       "test_key",
		Value:     "test_value",
		Tokens:    []string{"token1", "token2"},
		Policy:    [][]string{{"policy1"}, {"policy2"}},
		PolicyMap: map[string]int{"key1": 1, "key2": 2},
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ast.copy()
	}
}
