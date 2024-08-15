package utils

import (
	"slices"

	"dnevnik-rg.ru/internal/models"
)

func ValidateHistoryClasses(source []models.ShortStringClassInfo, filteredCoachID, filteredPupilID *int) (*[]models.ShortStringClassInfo, error) {
	if *filteredCoachID == 0 && *filteredPupilID == 0 {
		return &source, nil
	}

	validated := make([]models.ShortStringClassInfo, 0, len(source))

	for idx := range source {
		if *filteredPupilID != 0 && slices.Contains(
			source[idx].PupilsKeys,
			*filteredPupilID,
		) && *filteredCoachID == 0 {
			validated = append(validated, source[idx])
			continue
		}

		if *filteredCoachID != 0 && *filteredPupilID == 0 &&
			source[idx].CoachKey == *filteredCoachID {
			validated = append(validated, source[idx])
			continue
		}

		if *filteredCoachID != 0 &&
			*filteredPupilID != 0 &&
			slices.Contains(
				source[idx].PupilsKeys,
				*filteredPupilID,
			) && source[idx].CoachKey == *filteredCoachID {
			validated = append(validated, source[idx])
			continue
		}
	}

	return &validated, nil
}
