package photoprism

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/photoprism/photoprism/internal/ai/face"
	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/pkg/rnd"
)

// TestFaces_CorrectionCollision pins that correcting one face in a named cluster to another existing
// person records a collision that later runs keep, and moves the face to a face of that person.
func TestFaces_CorrectionCollision(t *testing.T) {
	w := isolatedTestFaces(t, "facescorrection")

	carol := consensusTestSubject(t, "Correction Carol")
	dave := consensusTestSubject(t, "Correction Dave")

	f := entity.NewFace(carol.SubjUID, entity.SrcAuto, face.Embeddings{face.FixtureEmbedding(7301)}, face.EmbeddingModelName())
	require.NotNil(t, f)
	require.NoError(t, f.Create())
	require.NoError(t, f.Updates(entity.Values{"samples": 5}))

	newMarker := func(t *testing.T, fraction float64, seed uint64) string {
		t.Helper()

		dist := fraction * f.AcceptDist()
		emb := face.Embeddings{face.FixtureEmbeddingAt(f.Embedding(), dist, seed)}
		m := entity.Marker{
			MarkerUID:      rnd.GenerateUID('m'),
			FileUID:        consensusTestFileUID,
			MarkerType:     entity.MarkerFace,
			MarkerSrc:      entity.SrcImage,
			SubjUID:        carol.SubjUID,
			SubjSrc:        entity.SrcAuto,
			FaceID:         f.ID,
			FaceDist:       dist,
			EmbeddingsJSON: emb.JSON(),
			EmbedModel:     f.EmbedModel,
			Size:           face.ClusterSizeThreshold,
			Score:          face.ClusterScore("") + 10,
			MatchedAt:      entity.TimeStamp(),
			W:              0.1,
			H:              0.1,
		}

		require.NoError(t, entity.UnscopedDb().Create(&m).Error)

		return m.MarkerUID
	}

	near := []string{newMarker(t, 0.2, 1), newMarker(t, 0.2, 2)}
	corrected := newMarker(t, 0.6, 3)

	m := entity.FindMarker(corrected)
	require.NotNil(t, m)
	changed, err := m.SetName(dave.SubjName, entity.SrcManual)
	require.NoError(t, err)
	require.True(t, changed)
	require.NoError(t, m.Save())

	check := func(t *testing.T) {
		t.Helper()

		cluster := entity.FindFace(f.ID)
		require.NotNil(t, cluster, "the cluster is kept")
		assert.Equal(t, carol.SubjUID, cluster.SubjUID)
		assert.Equal(t, 1, cluster.Collisions, "one correction is one collision")

		got := entity.FindMarker(corrected)
		require.NotNil(t, got)
		assert.Equal(t, dave.SubjUID, got.SubjUID)
		assert.NotEqual(t, f.ID, got.FaceID, "the corrected face does not return to the cluster")

		if own := entity.FindFace(got.FaceID); assert.NotNil(t, own, "it has a face of its own person") {
			assert.Equal(t, dave.SubjUID, own.SubjUID)
		}

		for _, uid := range near {
			assert.Equal(t, carol.SubjUID, entity.FindMarker(uid).SubjUID)
		}

		s := entity.FindSubject(carol.SubjUID)
		require.NotNil(t, s)
		assert.False(t, s.Deleted())
	}

	t.Run("Corrected", check)
	t.Run("Run", func(t *testing.T) {
		require.NoError(t, w.Start(FacesOptions{}))
		check(t)
	})
	t.Run("ForcedRun", func(t *testing.T) {
		require.NoError(t, w.Start(FacesOptions{Force: true}))
		check(t)
	})
}
