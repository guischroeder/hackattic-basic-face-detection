(ns hackattic-basic-face-detection.challenge
  (:require [hackattic-basic-face-detection.api.hackattic :as api]
            [hackattic-basic-face-detection.detection.face :as face]))

(defn solve-face-detection
  "Orchestrates the workflow for solving the face detection challenge"
  [access-token]
  (println "Fetching problem...")
  (let [problem (api/get-problem access-token)
        image-url (get problem :image_url)

        _ (println "Downloading image from" image-url)
        image-path (api/download-image image-url "resources/challenge-image.jpg")

        _ (println "Detecting faces...")
        face-coordinates (face/detect-faces image-path)

        _ (println "Submitting solution with" (count face-coordinates) "faces...")
        result (api/submit-solution access-token face-coordinates)]

    (println "Result from Hackattic:" result)
    result))
