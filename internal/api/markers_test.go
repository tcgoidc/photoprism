package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/photoprism/photoprism/internal/ai/face"
	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/form"
	"github.com/photoprism/photoprism/pkg/authn"
	"github.com/photoprism/photoprism/pkg/rnd"
)

func TestCreateMarker(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		app, router, _ := NewApiTest()

		GetPhoto(router)
		CreateMarker(router)

		r := PerformRequest(app, "GET", "/api/v1/photos/ps6sg6be2lvl0y11")

		assert.Equal(t, http.StatusOK, r.Code)

		photoUid := gjson.Get(r.Body.String(), "UID").String()
		fileUid := gjson.Get(r.Body.String(), "Files.0.UID").String()
		markerUid := gjson.Get(r.Body.String(), "Files.0.Markers.0.UID").String()

		assert.NotEmpty(t, photoUid)
		assert.NotEmpty(t, fileUid)
		assert.NotEmpty(t, markerUid)

		u := "/api/v1/markers"

		frm := form.Marker{
			FileUID:       fileUid,
			MarkerType:    "face",
			X:             0.303519,
			Y:             0.260742,
			W:             0.548387,
			H:             0.365234,
			SubjSrc:       "",
			MarkerName:    "",
			MarkerReview:  false,
			MarkerInvalid: false,
		}

		if b, err := json.Marshal(frm); err != nil {
			t.Fatal(err)
		} else {
			t.Logf("POST %s", u)
			r = PerformRequestWithBody(app, "POST", u, string(b))
		}

		assert.Equal(t, http.StatusCreated, r.Code)
		newUID := gjson.Get(r.Body.String(), "UID").String()
		assert.NotEmpty(t, newUID)
		assert.Equal(t, "/api/v1/markers/"+newUID, r.Header().Get("Location"))
	})
	t.Run("SuccessWithName", func(t *testing.T) {
		app, router, _ := NewApiTest()

		GetPhoto(router)
		CreateMarker(router)

		r := PerformRequest(app, "GET", "/api/v1/photos/ps6sg6be2lvl0y11")

		assert.Equal(t, http.StatusOK, r.Code)

		photoUid := gjson.Get(r.Body.String(), "UID").String()
		fileUid := gjson.Get(r.Body.String(), "Files.0.UID").String()
		markerUid := gjson.Get(r.Body.String(), "Files.0.Markers.0.UID").String()

		assert.NotEmpty(t, photoUid)
		assert.NotEmpty(t, fileUid)
		assert.NotEmpty(t, markerUid)

		u := "/api/v1/markers"

		frm := form.Marker{
			FileUID:       fileUid,
			MarkerType:    "face",
			X:             0.303519,
			Y:             0.260742,
			W:             0.548387,
			H:             0.365234,
			SubjSrc:       "manual",
			MarkerName:    "Jens Mander",
			MarkerReview:  false,
			MarkerInvalid: false,
		}

		if b, err := json.Marshal(frm); err != nil {
			t.Fatal(err)
		} else {
			t.Logf("POST %s", u)
			r = PerformRequestWithBody(app, "POST", u, string(b))
		}

		assert.Equal(t, http.StatusCreated, r.Code)
		newUID := gjson.Get(r.Body.String(), "UID").String()
		assert.NotEmpty(t, newUID)
		assert.Equal(t, "/api/v1/markers/"+newUID, r.Header().Get("Location"))
		assert.Equal(t, "Jens Mander", gjson.Get(r.Body.String(), "Name").String())
		assert.Equal(t, "manual", gjson.Get(r.Body.String(), "SubjSrc").String())
	})
	t.Run("InvalidArea", func(t *testing.T) {
		app, router, _ := NewApiTest()

		GetPhoto(router)
		CreateMarker(router)

		r := PerformRequest(app, "GET", "/api/v1/photos/ps6sg6be2lvl0y11")

		assert.Equal(t, http.StatusOK, r.Code)

		photoUid := gjson.Get(r.Body.String(), "UID").String()
		fileUid := gjson.Get(r.Body.String(), "Files.0.UID").String()
		markerUid := gjson.Get(r.Body.String(), "Files.0.Markers.0.UID").String()

		assert.NotEmpty(t, photoUid)
		assert.NotEmpty(t, fileUid)
		assert.NotEmpty(t, markerUid)

		u := "/api/v1/markers"

		frm := form.Marker{
			FileUID:       fileUid,
			MarkerType:    "face",
			X:             0.5,
			Y:             0.5,
			W:             0,
			H:             0,
			SubjSrc:       "",
			MarkerName:    "",
			MarkerReview:  false,
			MarkerInvalid: false,
		}

		if b, err := json.Marshal(frm); err != nil {
			t.Fatal(err)
		} else {
			t.Logf("POST %s", u)
			r = PerformRequestWithBody(app, "POST", u, string(b))
		}

		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
}

func TestUpdateMarker(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		app, router, _ := NewApiTest()

		GetPhoto(router)
		UpdateMarker(router)

		r := PerformRequest(app, "GET", "/api/v1/photos/ps6sg6be2lvl0y11")

		assert.Equal(t, http.StatusOK, r.Code)

		photoUid := gjson.Get(r.Body.String(), "UID").String()
		fileUid := gjson.Get(r.Body.String(), "Files.0.UID").String()
		markerUid := gjson.Get(r.Body.String(), "Files.0.Markers.0.UID").String()

		assert.NotEmpty(t, photoUid)
		assert.NotEmpty(t, fileUid)
		assert.NotEmpty(t, markerUid)

		u := fmt.Sprintf("/api/v1/markers/%s", markerUid)

		var m = form.Marker{
			SubjSrc:       "manual",
			MarkerInvalid: true,
			MarkerName:    "Foo",
		}

		if b, err := json.Marshal(m); err != nil {
			t.Fatal(err)
		} else {
			t.Logf("PUT %s", u)
			r = PerformRequestWithBody(app, "PUT", u, string(b))
		}

		assert.Equal(t, http.StatusOK, r.Code)
	})
	t.Run("NonPrimaryFile", func(t *testing.T) {
		app, router, _ := NewApiTest()

		UpdateMarker(router)

		r := PerformRequestWithBody(app, "PUT", "/api/v1/markers/ms6sg6b1wowu1000", "test")

		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
	t.Run("BadRequestFileAndPhotouidNotMatching", func(t *testing.T) {
		app, router, _ := NewApiTest()

		UpdateMarker(router)

		r := PerformRequestWithBody(app, "PUT", "/api/v1/markers/ms6sg6b1wowu1000", "test")

		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
	t.Run("FileNotExisting", func(t *testing.T) {
		app, router, _ := NewApiTest()

		UpdateMarker(router)

		r := PerformRequestWithBody(app, "PUT", "/api/v1/markers/1112", "test")

		assert.Equal(t, http.StatusNotFound, r.Code)
	})
	t.Run("MarkerNotExisting", func(t *testing.T) {
		app, router, _ := NewApiTest()

		UpdateMarker(router)

		r := PerformRequestWithBody(app, "PUT", "/api/v1/markers/1112", "test")

		assert.Equal(t, http.StatusNotFound, r.Code)
	})
	t.Run("EmptyPhotouid", func(t *testing.T) {
		app, router, _ := NewApiTest()

		UpdateMarker(router)

		r := PerformRequestWithBody(app, "PUT", "/api/v1/markers/ms6sg6b1wowu1000", "test")

		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
	t.Run("UpdateClusterWithExistingSubject", func(t *testing.T) {
		app, router, _ := NewApiTest()

		UpdateMarker(router)

		var m = form.Marker{
			SubjSrc:       "manual",
			MarkerInvalid: false,
			MarkerName:    "Actress A",
		}

		if b, err := json.Marshal(m); err != nil {
			t.Fatal(err)
		} else {
			r := PerformRequestWithBody(app, "PUT", "/api/v1/markers/ms6sg6b1wowuy666", string(b))

			assert.Equal(t, http.StatusOK, r.Code)

			ClearMarkerSubject(router)

			r = PerformRequestWithBody(app, "DELETE", "/api/v1/markers/ms6sg6b1wowuy666/subject", "")

			assert.Equal(t, http.StatusOK, r.Code)
		}
	})
	t.Run("UpdateClusterWithExistingSubjectTwo", func(t *testing.T) {
		app, router, _ := NewApiTest()

		UpdateMarker(router)

		var m = form.Marker{
			SubjSrc:       "manual",
			MarkerInvalid: false,
			MarkerName:    "Actress A",
		}

		if b, err := json.Marshal(m); err != nil {
			t.Fatal(err)
		} else {
			r := PerformRequestWithBody(app, "PUT", "/api/v1/markers/ms6sg6b1wowuy666", string(b))

			assert.Equal(t, http.StatusOK, r.Code)

			ClearMarkerSubject(router)

			r = PerformRequestWithBody(app, "DELETE", "/api/v1/markers/ms6sg6b1wowuy666/subject", "")

			assert.Equal(t, http.StatusOK, r.Code)
		}
	})
	t.Run("InvalidBody", func(t *testing.T) {
		app, router, _ := NewApiTest()

		UpdateMarker(router)

		var m = struct {
			ID      int
			Type    string
			Src     int
			Name    int
			SubjUID string
			SubjSrc string
			FaceID  string
		}{ID: 8,
			Type:    "face",
			Src:     123,
			Name:    456,
			SubjUID: "js6sg6b1h1njaaac",
			SubjSrc: "manual",
			FaceID:  "GMH5NISEEULNJL6RATITOA3TMZXMTMCI"}
		if b, err := json.Marshal(m); err != nil {
			t.Fatal(err)
		} else {
			r := PerformRequestWithBody(app, "PUT", "/api/v1/markers/ms6sg6b1wowuy666", string(b))

			assert.Equal(t, http.StatusBadRequest, r.Code)
		}
	})
}

func TestClearMarkerSubject(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		app, router, _ := NewApiTest()

		GetPhoto(router)
		ClearMarkerSubject(router)

		photoResp := PerformRequest(app, "GET", "/api/v1/photos/ps6sg6be2lvl0y11")

		if photoResp == nil {
			t.Fatal("response is nil")
		}

		assert.Equal(t, http.StatusOK, photoResp.Code)

		if photoResp.Body.String() == "" {
			t.Fatal("body is empty")
		}

		photoUid := gjson.Get(photoResp.Body.String(), "UID").String()
		fileUid := gjson.Get(photoResp.Body.String(), "Files.0.UID").String()
		markerUid := gjson.Get(photoResp.Body.String(), "Files.0.Markers.0.UID").String()

		assert.NotEmpty(t, photoUid)
		assert.NotEmpty(t, fileUid)
		assert.NotEmpty(t, markerUid)

		u := fmt.Sprintf("/api/v1/markers/%s/subject", markerUid)

		// t.Logf("DELETE %s", u)

		resp := PerformRequestWithBody(app, "DELETE", u, "")

		assert.Equal(t, http.StatusOK, resp.Code)
	})
	t.Run("NonPrimaryFile", func(t *testing.T) {
		app, router, _ := NewApiTest()

		ClearMarkerSubject(router)

		r := PerformRequestWithBody(app, "DELETE", "/api/v1/markers/ms6sg6b1wowu1000/subject", "")

		assert.Equal(t, http.StatusOK, r.Code)
	})
}

// TestUpdateMarker_NamedCluster pins that naming one marker in a cluster named after another
// person changes only that marker.
func TestUpdateMarker_NamedCluster(t *testing.T) {
	app, router, conf := NewApiTest()
	UpdateMarker(router)

	prevAuthMode := conf.AuthMode()
	conf.SetAuthMode(config.AuthModePasswd)
	t.Cleanup(func() { conf.SetAuthMode(prevAuthMode) })

	sess := entity.NewSession(conf.SessionMaxAge(), 0)
	sess.SetUser(entity.UserFixtures.Pointer("alice"))
	sess.SetScope("files")
	sess.SetProvider(authn.ProviderApplication)
	require.NoError(t, sess.Create())
	t.Cleanup(func() { entity.UnscopedDb().Unscoped().Delete(sess) })
	require.False(t, sess.SeesPrivatePeople())

	person := entity.NewSubject("Named Cluster Person", entity.SubjPerson, entity.SrcManual)
	require.NotNil(t, person)
	require.NoError(t, person.Create())
	t.Cleanup(func() { entity.UnscopedDb().Delete(&entity.Subject{}, "subj_uid = ?", person.SubjUID) })
	markPrivate(t, person, false)

	f := entity.NewFace(person.SubjUID, entity.SrcAuto, face.Embeddings{face.FixtureEmbedding(7201)}, face.EmbeddingModelName())
	require.NotNil(t, f)
	require.NoError(t, f.Create())
	t.Cleanup(func() { entity.UnscopedDb().Delete(&entity.Face{}, "id = ?", f.ID) })

	newMarker := func(subjUID, subjSrc string) string {
		m := entity.Marker{
			MarkerUID:  rnd.GenerateUID('m'),
			FileUID:    entity.FileFixtures.Get("exampleDNGFile.dng").FileUID,
			MarkerType: entity.MarkerFace,
			SubjUID:    subjUID,
			SubjSrc:    subjSrc,
			FaceID:     f.ID,
			FaceDist:   0.1,
			EmbedModel: f.EmbedModel,
			MatchedAt:  entity.TimeStamp(),
			W:          0.1,
			H:          0.1,
		}

		require.NoError(t, entity.UnscopedDb().Create(&m).Error)
		t.Cleanup(func() { entity.UnscopedDb().Delete(&entity.Marker{}, "marker_uid = ?", m.MarkerUID) })

		return m.MarkerUID
	}

	auto := newMarker(person.SubjUID, entity.SrcAuto)
	unnamed := newMarker("", entity.SrcAuto)
	rejected := newMarker("", entity.SrcManual)

	b, err := json.Marshal(form.Marker{SubjSrc: entity.SrcManual, MarkerName: "Named Cluster Other"})
	require.NoError(t, err)

	r := AuthenticatedRequestWithBody(app, http.MethodPut, "/api/v1/markers/"+rejected, string(b), sess.AuthToken())
	require.Equal(t, http.StatusOK, r.Code, r.Body.String())

	t.Cleanup(func() { entity.UnscopedDb().Delete(&entity.Subject{}, "subj_name = ?", "Named Cluster Other") })

	named := entity.FindMarker(rejected)
	require.NotNil(t, named)
	require.NotEmpty(t, named.SubjUID, "the named marker changes")
	assert.NotEqual(t, person.SubjUID, named.SubjUID)

	assert.Equal(t, person.SubjUID, entity.FindMarker(auto).SubjUID)
	assert.Empty(t, entity.FindMarker(unnamed).SubjUID)
	assert.Equal(t, person.SubjUID, entity.FindFace(f.ID).SubjUID)
}
