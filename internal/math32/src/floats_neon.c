#include <arm_neon.h>

void _dot_product_neon(float *a, float *b, long n, float* result)
{
    int epoch = n / 8; // Number of full vectors (size 4) to process
    int remain = n % 8; // Number of elements left after processing full vectors
    
    float32x4_t sum1 = vdupq_n_f32(0.0f);
    float32x4_t sum2 = vdupq_n_f32(0.0f);
    float32x4_t v1, v2;
    
    // Vectorized computation loop with loop unrolling
    for (int i = 0; i < epoch; ++i)
    {
        v1 = vld1q_f32(a);
        v2 = vld1q_f32(b);
        sum1 = vmlaq_f32(sum1, v1, v2);
        a += 4;
        b += 4;

        v1 = vld1q_f32(a);
        v2 = vld1q_f32(b);
        sum2 = vmlaq_f32(sum2, v1, v2);
        a += 4;
        b += 4;
    }

    // Process remaining elements
    float remain_sum = 0.0f;
    for (int i = 0; i < remain; ++i)
    {
        remain_sum += a[i] * b[i];
    }

    // add remaining sum to sum1[0]
    sum1 = vsetq_lane_f32(vgetq_lane_f32(sum1, 0) + remain_sum, sum1, 0);
    
    // Horizontal sum of the vectors
    sum1 = vaddq_f32(sum1, sum2);
    sum1 = vpaddq_f32(sum1, sum1);
    sum1 = vpaddq_f32(sum1, sum1);
    
    // Extract the final result
    float32x2_t sum2_lanes = vget_low_f32(sum1);
    *result = vget_lane_f32(sum2_lanes, 0);
}

void _squared_l2_neon(float *a, float *b, long n, float *result)
{
    float32x4_t sumQuad = vdupq_n_f32(0.0f);
    long quadCount = n / 4;
    long remainder = n % 4;

    while (quadCount--) {
        float32x4_t aQuad = vld1q_f32(a);
        float32x4_t bQuad = vld1q_f32(b);
        float32x4_t diffQuad = vsubq_f32(aQuad, bQuad);
        sumQuad = vmlaq_f32(sumQuad, diffQuad, diffQuad);
        a += 4;
        b += 4;
    }

    // Now, after bulk operation, perform operation on remaining vars
    float32x2_t sumDouble = vadd_f32(vget_low_f32(sumQuad), vget_high_f32(sumQuad));

    float32_t sum = vget_lane_f32(sumDouble, 0) + vget_lane_f32(sumDouble, 1);

    // Scalar operation on remaining vars
    for (; remainder; --remainder) {
        float diff = *a++ - *b++;
        sum += diff * diff;
    }

    *result = sum;
}
