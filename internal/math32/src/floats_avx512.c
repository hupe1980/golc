#include <immintrin.h>

void _dot_product_avx512(float* vec1, float* vec2, long n, float* result)
{
    __m512 res = _mm512_setzero_ps();
    long vectorized_end = n - (n % 16);

    for (long i = 0; i < vectorized_end; i += 16)
    {
        __m512 float1 = _mm512_loadu_ps(vec1 + i);
        __m512 float2 = _mm512_loadu_ps(vec2 + i);
        __m512 mult = _mm512_mul_ps(float1, float2);
        res = _mm512_add_ps(res, mult);
    }

    float vec[16] __attribute__((aligned(64))) = {0};
    _mm512_store_ps(vec, res);

    *result = 0;
    for (int i = 0; i<16; i++)
    {
        *result+= vec[i];
    }

    // handle the leftover part
    for (long i = vectorized_end; i < n; ++i)
    {
        *result += vec1[i] * vec2[i];
    }
}

void _squared_l2_avx512(float* vec1, float* vec2, long n, float* result)
{
    __m512 sum = _mm512_setzero_ps();

    long i;
    for(i=0; i<n/16; i++)
    {
        __m512 v1 = _mm512_loadu_ps(&vec1[i*16]);
        __m512 v2 = _mm512_loadu_ps(&vec2[i*16]);
        __m512 diff = _mm512_sub_ps(v1, v2);
        __m512 sq = _mm512_mul_ps(diff, diff);
        sum = _mm512_add_ps(sum, sq);
    }

    float temp[16];
    _mm512_storeu_ps(temp, sum);
    float total = temp[0] + temp[1] + temp[2] + temp[3] + temp[4] + temp[5] + temp[6] + temp[7]
            + temp[8] + temp[9] + temp[10] + temp[11] + temp[12] + temp[13] + temp[14] + temp[15];

    // process the leftovers
    for(long j=i*16; j<n; j++)
    {
        float diff = vec1[j] - vec2[j];
        total += diff * diff;
    }

    *result = total;
}
