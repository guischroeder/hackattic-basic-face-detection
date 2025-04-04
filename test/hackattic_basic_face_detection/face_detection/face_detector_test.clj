(ns hackattic-basic-face-detection.face-detection.face-detector-test
  (:require [clojure.test :refer [deftest is testing]]
            [hackattic-basic-face-detection.face-detection.face-detector :as face])
  (:import [org.bytedeco.opencv.opencv_core Mat RectVector Rect]))

(deftest detection-params-test
  (testing "Detection parameters have expected values"
    (is (> (:scale-factor face/detection-params) 1.0) "Scale factor should be greater than 1.0")
    (is (>= (:min-neighbors face/detection-params) 0) "Min neighbors should be non-negative")
    (is (pos? (:min-face-size face/detection-params)) "Min face size should be positive")))

(deftest convert-to-coordinates-test
  (testing "Converting RectVector to coordinates works correctly"
    (let [rect (doto (Rect.) (.x 10) (.y 20) (.width 30) (.height 40))
          rect-vector (doto (RectVector.) (.push_back rect))
          result (face/convert-to-coordinates rect-vector)]

      (is (vector? result) "Result should be a vector")
      (is (= 1 (count result)) "Should have one face")
      (is (= {:x 10 :y 20 :width 30 :height 40} (first result)) "Coordinates should match"))))

(deftest detect-faces-test
  (testing "Face detection with mocked dependencies"
    (let [test-rect (doto (Rect.) (.x 100) (.y 200) (.width 50) (.height 60))
          test-faces (doto (RectVector.) (.push_back test-rect))]

      (with-redefs [face/ensure-cascade-file (fn [] nil)
                    face/load-and-preprocess-image (fn [_] (Mat.))
                    face/detect-faces-in-image (fn [_] test-faces)]

        (let [result (face/detect-faces "test-image.jpg")]
          (is (vector? result) "Result should be a vector")
          (is (= 1 (count result)) "Should detect one face")
          (is (= {:x 100 :y 200 :width 50 :height 60} (first result)) "Face coordinates should match"))))))
