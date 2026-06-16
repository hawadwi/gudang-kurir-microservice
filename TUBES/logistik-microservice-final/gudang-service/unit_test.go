package main

import (
	"testing"
)

func TestStartSortingSuccess(t *testing.T) {
	repo := &PackageRepository{db: nil}
	service := NewSortingService(repo)

	pkg := &Package{
		UserID:        123,
		Resi:          "RES001",
		NamaBarang:    "Laptop",
		Berat:         2,
		WarehouseZone: "Jakarta",
		Status:        "pending",
	}

	err := service.StartSorting(pkg)
	if err != nil {
		t.Errorf("StartSorting failed: %v", err)
	}

	if pkg.Status != "sorting" {
		t.Errorf("Expected status 'sorting', got '%s'", pkg.Status)
	}
}

func TestStartSortingNilResi(t *testing.T) {
	repo := &PackageRepository{db: nil}
	service := NewSortingService(repo)

	pkg := &Package{
		Resi:          "",
		WarehouseZone: "Jakarta",
		Status:        "pending",
	}

	err := service.StartSorting(pkg)
	if err == nil {
		t.Error("Expected error for empty resi")
	}
}

func TestCompleteSortingSuccess(t *testing.T) {
	repo := &PackageRepository{db: nil}
	service := NewSortingService(repo)

	pkg := &Package{
		Resi:          "RES001",
		Status:        "sorting",
		WarehouseZone: "Jakarta",
	}

	err := service.CompleteSorting(pkg)
	if err != nil {
		t.Errorf("CompleteSorting failed: %v", err)
	}

	if pkg.Status != "ready" {
		t.Errorf("Expected status 'ready', got '%s'", pkg.Status)
	}

	if pkg.SortedAt == nil {
		t.Error("SortedAt should not be nil")
	}
}

func TestCompleteSortingNil(t *testing.T) {
	repo := &PackageRepository{db: nil}
	service := NewSortingService(repo)

	err := service.CompleteSorting(nil)
	if err == nil {
		t.Error("Expected error for nil package")
	}
}

func TestCompleteSortingInvalidStatus(t *testing.T) {
	repo := &PackageRepository{db: nil}
	service := NewSortingService(repo)

	pkg := &Package{
		Resi:          "RES001",
		Status:        "pending",
		WarehouseZone: "Jakarta",
	}

	err := service.CompleteSorting(pkg)
	if err == nil {
		t.Error("Expected error for non-sorting status")
	}
}

func BenchmarkStartSorting(b *testing.B) {
	repo := &PackageRepository{db: nil}
	service := NewSortingService(repo)

	pkg := &Package{
		UserID:        123,
		Resi:          "RES001",
		WarehouseZone: "Jakarta",
		Status:        "pending",
	}

	for i := 0; i < b.N; i++ {
		service.StartSorting(pkg)
		pkg.Status = "pending"
	}
}
