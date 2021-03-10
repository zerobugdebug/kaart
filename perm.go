package main

func nextPerm(p []int) {
	for i := len(p) - 1; i >= 0; i-- {
		if i == 0 || p[i] < len(p)-i-1 {
			p[i]++
			return
		}
		p[i] = 0
	}
}

func getPerm(orig, p []int) []int {
	result := append([]int{}, orig...)
	for i, v := range p {
		result[i], result[i+v] = result[i+v], result[i]
	}
	return result
}

//GetAllPermutations will generate all possible permutations of specific int slice
func GetAllPermutations(orig []int) [][]int {
	var permutations [][]int
	for p := make([]int, len(orig)); p[0] < len(p); nextPerm(p) {
		permutations = append(permutations, getPerm(orig, p))
	}
	return permutations
}

//GetAllPermutationsForSum will generate all possible permutations for k elements, where sum of all elements equal to totalSum
func GetAllPermutationsForSum(k, totalSum int) [][]int {
	var result [][]int
	var output []int
	getPermutationsForSum(k-1, 0, totalSum, output, &result)
	//Logger.Debug(output)
	//Logger.Debug(result)
	return result

}

func getPermutationsForSum(k, currentSum, totalSum int, output []int, result *[][]int) {
	if k == 0 {
		output = append(output, totalSum-currentSum)
		outputCopy := make([]int, len(output))
		copy(outputCopy, output)
		*result = append(*result, outputCopy)
	} else {

		for i := 0; i <= totalSum-currentSum; i++ {
			output = append(output, i)
			//Logger.Debug(k, currentSum, totalSum, output, result)
			getPermutationsForSum(k-1, currentSum+i, totalSum, output, result)
			//Logger.Debug(k, currentSum, totalSum, output, result)
			output = output[:len(output)-1]
		}
	}
}
