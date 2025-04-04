(ns hackattic-basic-face-detection.challenge-test
  (:require [clojure.test :refer [deftest is testing]]
            [hackattic-basic-face-detection.challenge :as challenge]
            [hackattic-basic-face-detection.api.hackattic :as api]
            [hackattic-basic-face-detection.detection.face :as face]))

(deftest solve-face-detection-test
  (testing "End-to-end face detection challenge workflow"
    (let [test-token "test-token"
          test-image-url "https://example.com/test-image.jpg"
          test-image-path "resources/challenge-image.jpg"
          test-face-coordinates [{:x 100 :y 200 :width 50 :height 60}]
          test-api-result {:result "success"}]

      (with-redefs [api/get-problem (fn [token]
                                      (is (= test-token token))
                                      L                      {:image_url test-image-url})

                    api/download-image (fn [url path]
                                         (is (= test-image-url url))
                                         (is (= test-image-path path))
                                         test-image-path)

                    api/submit-solution (fn [token coordinates]
                                          (is (= test-token token))
                                          (is (= test-face-coordinates coordinates))
                                          test-api-result)

                    face/detect-faces (fn [path]
                                        (is (= test-image-path path))
                                        test-face-coordinates)]

        (let [result (challenge/solve-face-detection test-token)]
          (is (= test-api-result result))
          (is (= "success" (:result result))))))))
