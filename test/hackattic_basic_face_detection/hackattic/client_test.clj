(ns hackattic-basic-face-detection.hackattic.client-test
  (:require [clojure.test :refer [deftest is testing]]
            [hackattic-basic-face-detection.hackattic.client :as client]
            [clj-http.client :as http]
            [cheshire.core :as json]
            [clojure.java.io :as io]))

(deftest get-problem-test
  (testing "Getting problem from API"
    (let [test-token "test-token"
          mock-response {:status 200
                         :body {:image_url "https://example.com/test-image.jpg"}}]

      (with-redefs [http/get (fn [url _]
                               (is (= (str client/base-url "/problem?access_token=" test-token) url))
                               mock-response)]

        (let [result (client/get-problem test-token)]
          (is (= (:body mock-response) result))
          (is (= "https://example.com/test-image.jpg" (:image_url result))))))))

(deftest download-image-test
  (testing "Downloading image from URL"
    (let [test-url "https://example.com/test-image.jpg"
          test-path "test-output.jpg"
          mock-input-stream (java.io.ByteArrayInputStream. (.getBytes "test-image-data"))
          mock-output-stream (java.io.ByteArrayOutputStream.)]

      (with-redefs [io/make-parents (fn [_] nil)
                    io/input-stream (fn [url]
                                      (is (= test-url url))
                                      mock-input-stream)
                    io/output-stream (fn [path]
                                       (is (= test-path path))
                                       mock-output-stream)]

        (let [result (client/download-image test-url test-path)]
          (is (= test-path result)))))))

(deftest submit-solution-test
  (testing "Submitting solution to API"
    (let [test-token "test-token"
          test-coordinates [{:x 10 :y 20 :width 30 :height 40}]
          expected-solution {:faces test-coordinates}
          mock-response {:status 200
                         :body {:result "success"}}]

      (with-redefs [http/post (fn [url {:keys [body content-type]}]
                                (is (= (str client/base-url "/solve?access_token=" test-token) url))
                                (is (= :json content-type))
                                (is (= (json/generate-string expected-solution) body))
                                mock-response)]

        (let [result (client/submit-solution test-token test-coordinates)]
          (is (= (:body mock-response) result))
          (is (= "success" (:result result))))))))
