#include <immintrin.h> // AVX intrinsics

void _dot_product_avx(float *a, float *b, long n, float *res)
{
    __m256 sumvec = _mm256_setzero_ps();
    long i;
    for (i = 0; i < n - (n % 8); i += 8)
    {
        __m256 avec = _mm256_loadu_ps(&a[i]);
        __m256 bvec = _mm256_loadu_ps(&b[i]);
        __m256 mult = _mm256_mul_ps(avec, bvec);
        sumvec = _mm256_add_ps(sumvec, mult);
    }

    // handle leftovers if n is not a multiple of 8
    float temp[8] = {0.0};
    _mm256_storeu_ps(temp, sumvec);
    for (int j = 0; j < 8; j++)
    {
        *res += temp[j];
    }
    
    for (; i < n; i++)
    {
        *res += a[i] * b[i];
    }
}

void _squared_l2_avx(float *vec1, float *vec2, long n, float *result)
{
    __m256 res = _mm256_setzero_ps(); // initialize result to zero
    long i = 0;

    // iterate over the vectors in chunks of 8
    for (; i < n - (n % 8); i += 8)
    {
        __m256 vec1_8 = _mm256_loadu_ps(vec1 + i);
        __m256 vec2_8 = _mm256_loadu_ps(vec2 + i);

        // find the difference
        __m256 diff = _mm256_sub_ps(vec1_8, vec2_8);

        // square the difference
        __m256 sqrd_diff = _mm256_mul_ps(diff, diff);
        res = _mm256_add_ps(res, sqrd_diff);
    }

    // accumulate the residual elements using scalar operations
    float temp[8];
    _mm256_storeu_ps(temp, res);
    float sum = 0.0f;
    for (; i < n; i++)
    {
        float diff = vec1[i] - vec2[i];
        sum += diff * diff;
    }

    // add up the results from AVX and scalar operations
    for (int j = 0; j < 8; j++)
    {
        sum += temp[j];
    }

    // store the final result
    *result = sum;
}