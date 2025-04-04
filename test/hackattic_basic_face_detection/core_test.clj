(ns hackattic-basic-face-detection.core-test
  (:require [clojure.test :refer [deftest is testing]]
            [hackattic-basic-face-detection.core :as core]
            [hackattic-basic-face-detection.hackattic.client :as client]
            [hackattic-basic-face-detection.face_detection.face_detector :as face]))

(deftest solve-face-detection-test
  (testing "End-to-end face detection challenge workflow"
    (let [test-token "test-token"
          test-image-url "https://example.com/test-image.jpg"
          test-image-path "resources/challenge-image.jpg"
          test-face-coordinates [{:x 100 :y 200 :width 50 :height 60}]
          test-api-result {:result "success"}]

      (with-redefs [client/get-problem (fn [token]
                                      (is (= test-token token))
                                      {:image_url test-image-url})

                    client/download-image (fn [url path]
                                         (is (= test-image-url url))
                                         (is (= test-image-path path))
                                         test-image-path)

                    client/submit-solution (fn [token coordinates]
                                          (is (= test-token token))
                                          (is (= test-face-coordinates coordinates))
                                          test-api-result)

                    face/detect-faces (fn [path]
                                        (is (= test-image-path path))
                                        test-face-coordinates)]

        (let [result (core/solve-face-detection test-token)]
          (is (= test-api-result result))
          (is (= "success" (:result result))))))))

(deftest get-access-token-test
  (testing "Getting access token from args"
    (let [args ["test-token"]]
      (is (= "test-token" (core/get-access-token args)))))
  
  (testing "Getting access token from environment"
    (System/setProperty "HACKATTIC_TOKEN" "env-token")
    (is (= "env-token" (core/get-access-token []))))
  
  (testing "Throws exception when no token available"
    (System/clearProperty "HACKATTIC_TOKEN")
    (is (thrown? IllegalArgumentException (core/get-access-token [])))))
