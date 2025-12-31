import type { Ref } from "vue"

const FIXED_TIME = 3000;

export const createIncresingValue = (dynamicValue: Ref<number>, targetValue: number, step?: number) => {
  const totalSteps = Math.ceil(FIXED_TIME / 2);
  step = Math.max(Math.ceil((targetValue - dynamicValue.value) / totalSteps), 1);

  const timer = setInterval(() => {
    if (targetValue - dynamicValue.value < step) {
      dynamicValue.value = targetValue;
      clearInterval(timer);
    } else {
      dynamicValue.value += step;
    }
  }, 2);
}
