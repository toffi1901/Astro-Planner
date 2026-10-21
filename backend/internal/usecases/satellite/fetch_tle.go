package satellite

import (
	"astro-planner/backend/internal/entities"
	"astro-planner/backend/internal/repository"
	"bufio"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type FetchTLEUseCase struct {
	Repo *repository.SatelliteRepository
}

func NewFetchTLEUseCase(repo *repository.SatelliteRepository) *FetchTLEUseCase {
	return &FetchTLEUseCase{Repo: repo}
}

func (uc *FetchTLEUseCase) Execute(ctx context.Context) error {
	url := "https://r4uab.ru/satonline.txt"

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("fetch tle error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("r4uab returned status: %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	var lines []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	for i := 0; i < len(lines); i += 3 {
		if i+2 >= len(lines) {
			break
		}

		name := lines[i]
		line1 := lines[i+1]
		line2 := lines[i+2]

		if len(line1) < 7 || len(line2) < 7 {
			continue
		}

		noradStr := strings.TrimSpace(line1[2:7])
		noradID, err := strconv.ParseInt(noradStr, 10, 64)
		if err != nil {
			continue
		}

		// ФИЛЬТР: Оставляем только МКС (25544) и Хаббл (20580) просто, чтобы не было больно чичас это ВРЕМЕННО
		if noradID != 25544 && noradID != 20580 {
			continue
		}

		// простая классификация
		satType := "other"
		if noradID == 25544 {
			satType = "iss"
		} else if noradID == 20580 {
			satType = "scientific"
		}

		sat := &entities.Satellite{
			NoradID:        noradID,
			Name:           name,
			Type:           satType,
			TLELine1:       line1,
			TLELine2:       line2,
			LastUpdatedTLE: time.Now(),
			IsActive:       true,
		}

		if err := uc.Repo.Save(ctx, sat); err != nil {
			return fmt.Errorf("save satellite %d error: %w", noradID, err)
		}
	}

	fmt.Println("Успешно загружены TLE данные с R4UAB")
	return nil
}
