package utils

import (
	"fmt"
	"strings"
	"testing"
)

func TestOnlyCnes(t *testing.T) {
	// Usando uma struct mais idiomática com campos alinhados
	tests := []struct {
		name string
		input string
		want string
	}{
		// Casos principais de uso
		{"main_scenario", "1234567 - nome da unidade 123", "1234567"},
		{"with_hyphen", "1234567-Hospital Central", "1234567"},
		{"digits_only", "1234567", "1234567"},
		{"real_world_example", "2382857 - UBS Vila Nova São Paulo", "2382857"},
		
		// Casos com formatação diferente
		{"leading_spaces", "  1234567 - UBS Centro", "1234567"},
		{"leading_tabs", "\t1234567-UPA Norte", "1234567"},
		{"with_special_chars", "1234567@#$%^&*()_+ Unidade", "1234567"},
		{"with_newline", "1234567\nHospital Municipal", "1234567"},
		{"additional_numbers", "1234567 - Unidade 123 Setor 456", "1234567"},
		
		// Casos com mais de 7 dígitos
		{"more_than_7_digits", "12345678901234", "1234567"},
		
		// Casos especiais
		{"zeros", "0000000 - Unidade Teste", "0000000"},
		
		// Casos de erro/edge cases
		{"no_leading_numbers", "Hospital 1234567", ""},
		{"less_than_7_digits", "123456 - Posto de Saúde", ""},
		{"empty_string", "", ""},
		{"only_spaces", "   ", ""},
	}

	for _, tc := range tests { // tc (test case) é mais idiomático que tt
		t.Run(tc.name, func(t *testing.T) {
			got := OnlyCnes(tc.input)
			if got != tc.want {
				t.Errorf("OnlyCnes(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestOnlyCnes_Performance(t *testing.T) {
	// Testa comportamento com strings muito longas
	veryLongString := "1234567" + strings.Repeat(" very long string", 1000)
	got := OnlyCnes(veryLongString)
	want := "1234567"
	
	if got != want {
		t.Errorf("OnlyCnes with very long string: got %q, want %q", got, want)
	}
}

// BenchmarkOnlyCnes mede a performance da função OnlyCnes
func BenchmarkOnlyCnes(b *testing.B) {
	input := "1234567 - nome da unidade 123"
	
	b.ResetTimer() // Recomendado em benchmarks
	for i := 0; i < b.N; i++ {
		_ = OnlyCnes(input) // Atribuição blank para evitar otimizações
	}
}

// BenchmarkOnlyCnesLongString testa performance com strings longas
func BenchmarkOnlyCnesLongString(b *testing.B) {
	input := "1234567" + strings.Repeat(" long string", 100)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = OnlyCnes(input)
	}
}

// TestOnlyCnesVsOnlyNumber demonstra a diferença entre as funções
func TestOnlyCnesVsOnlyNumber(t *testing.T) {
	input := "1234567 - nome da unidade 123"
	
	cnesGot := OnlyCnes(input)
	numberGot := OnlyNumber(input)
	
	cnesWant := "1234567"          // Apenas primeiros 7 dígitos
	numberWant := "1234567123"     // Todos os números
	
	if cnesGot != cnesWant {
		t.Errorf("OnlyCnes(%q) = %q, want %q", input, cnesGot, cnesWant)
	}
	
	if numberGot != numberWant {
		t.Errorf("OnlyNumber(%q) = %q, want %q", input, numberGot, numberWant)
	}
	
	// Assegurar que comportamentos são diferentes
	if cnesGot == numberGot {
		t.Error("OnlyCnes and OnlyNumber should return different values")
	}
}

// Helper function para testes - boa prática para evitar repetição
func assertOnlyCnes(t *testing.T, input, want string) {
	t.Helper() // Marca como função auxiliar de teste
	got := OnlyCnes(input)
	if got != want {
		t.Errorf("OnlyCnes(%q) = %q, want %q", input, got, want)
	}
}

// ExampleOnlyCnes demonstra o uso da função OnlyCnes
// Exemplos são executados como testes e aparecem na documentação
func ExampleOnlyCnes() {
	// Caso de uso principal: extrair CNES de uma string com nome da unidade
	result := OnlyCnes("1234567 - Hospital Central")
	fmt.Println(result)
	// Output: 1234567
}

func ExampleOnlyCnes_withExtraNumbers() {
	// Demonstra que apenas os primeiros 7 dígitos são extraídos
	result := OnlyCnes("1234567 - Unidade 123 Setor 456")
	fmt.Println(result)
	// Output: 1234567
}

// TestOnlyCnes_TableDrivenWithHelper demonstra uso de função auxiliar
func TestOnlyCnes_TableDrivenWithHelper(t *testing.T) {
	// Casos rápidos usando helper function
	assertOnlyCnes(t, "1234567 - UBS", "1234567")
	assertOnlyCnes(t, "7654321-Hospital", "7654321")
	assertOnlyCnes(t, "0000000 Test", "0000000")
	assertOnlyCnes(t, "123456", "") // Menos de 7 dígitos
}
