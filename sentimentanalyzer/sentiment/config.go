package sentiment

// var plainAnalysisSystemPrompt string = `You are a markets expert.
// 	Analyze the sentiment of this financial news related to the given symbol and respond with one of the following words about the sentiment
// 	[positive, negative, neutral]. Respond with only one word.`
// var semanticAnalysisSystemPrompt string = `You are a markets expert.
// Perform semantic sentiment analysis of this financial news related to the given symbols and respond with one of the following words about the sentiment
// [positive, negative, neutral]. Respond with a well formatted json object where each key is a symbol and the value is the sentiment of that symbol in the news. Example: {"AAPL": "positive", "GOOG": "negative"}.`
// var classificationAnalysisSystemPrompt string
// var symbolCheckPrompt string = `You are a classification tool that knows about asset symbols and their names. Given a news headline and a symbol (in the format {headline};{symbol}) you will identify whether the headline has anything to do with the symbol or the company that the symbol represents. If it does answer with true, false otherwise. You can only answer with one word.`
