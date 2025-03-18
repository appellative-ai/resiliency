package limit

// Need to capture all traffic for a given time period.
//  1. Need to include rate limited requests
//  2. Need to have access to timeout value to determine when to start rate limiting
//  3. Determine saturation, percentile latency, and gradiant for a window
