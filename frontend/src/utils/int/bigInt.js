import Big from "big.js";

export const GetPercentScore = (score, percent) => {
    // const num1 = new BigNumber(score);
    // const num2 = new BigNumber(percent);
    let num1 = new Big(score);
    let num2 = new Big(percent);
  
    const result = num1.times(num2);
  
    return Number(result);
  }