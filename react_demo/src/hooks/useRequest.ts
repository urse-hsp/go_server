// hooks/useApi.ts
import useSWR from 'swr'
import type { SWRConfiguration } from 'swr'
import useSWRMutation, { type SWRMutationResponse } from 'swr/mutation'
import { fetcher } from '../utils/fetcher'
import useSWRInfinite, { type SWRInfiniteConfiguration } from 'swr/infinite'

const LIMIT_SIZE = 10
const LIMIT_NAME = 'limit'
const CURSOR_NAME = 'cursor'

// get
export function useGetRequest<T>(
  url: string | null, // url 为 null 时不发起请求
  options?: SWRConfiguration,
) {
  const { data, error, mutate, isValidating } = useSWR<T>(url, fetcher, {
    revalidateOnFocus: false,
    shouldRetryOnError: false,
    ...options,
  })

  return {
    data,
    loading: !error && !data && isValidating,
    error: error ? (error as Error).message : null,
    mutate,
  }
}

// post
export function useTriggerRequest<T>(url: string) {
  return useSWRMutation<T>(url, fetcher)
}

// get list
type InfiniteOptions<T> = {
  pageSize?: number
  getNextPageParam?: (lastPage: T, allPages: T[]) => any
  isEnd?: (lastPage: T) => boolean
  initialParams?: Record<string, any>
  swrOptions?: SWRInfiniteConfiguration
}
export function useGetRequestInfinite<T = any>(
  baseUrl: string,
  options?: InfiniteOptions<T>,
) {
  const {
    pageSize = LIMIT_SIZE,
    getNextPageParam,
    isEnd,
    initialParams = {},
    swrOptions,
  } = options || {}

  const getKey = (pageIndex: number, previousPageData: any) => {
    // 第一页
    if (pageIndex === 0) {
      const query = new URLSearchParams({
        ...initialParams,
        limit: String(pageSize),
      }).toString()
      return `${baseUrl}?${query}`
    }

    // 判断是否结束
    if (isEnd?.(previousPageData)) return null

    // cursor / page
    const nextParam = getNextPageParam
      ? getNextPageParam(previousPageData, [])
      : previousPageData?.nextCursor

    if (!nextParam) return null

    const query = new URLSearchParams({
      ...initialParams,
      [CURSOR_NAME]: String(nextParam),
      [LIMIT_NAME]: String(pageSize),
    }).toString()

    return `${baseUrl}?${query}`
  }

  const swr = useSWRInfinite<T>(getKey, fetcher, {
    revalidateFirstPage: false,
    ...swrOptions,
  })

  // 👉 平铺数据（重点）
  const list = swr.data ? swr.data.flatMap((item: any) => item.data || []) : []

  const isLoadingInitialData = !swr.data && !swr.error
  const isLoadingMore =
    isLoadingInitialData ||
    (swr.size > 0 && swr.data && typeof swr.data[swr.size - 1] === 'undefined')

  const isEmpty = swr.data?.[0]?.data?.length === 0
  const isReachingEnd =
    isEmpty || (swr.data && isEnd?.(swr.data[swr.data.length - 1]))

  return {
    ...swr,

    list, // 👈 扁平数组
    loading: isLoadingInitialData,
    loadingMore: isLoadingMore,
    noMore: isReachingEnd,

    loadMore: () => swr.setSize(swr.size + 1),
  }
}

// get list 最终版本
export function useGetRequestInfinite2<T = any>(
  baseUrl: string,
  options?: InfiniteOptions<T>,
) {
  const {
    pageSize = LIMIT_SIZE,
    getNextPageParam,
    isEnd,
    initialParams = {},
    swrOptions,
  } = options || {}

  const getKey = (pageIndex: number, previousPageData: any) => {
    if (pageIndex === 0) {
      const query = new URLSearchParams({
        ...initialParams,
        limit: String(pageSize),
      }).toString()
      return `${baseUrl}?${query}`
    }

    if (isEnd?.(previousPageData)) return null

    const nextParam = getNextPageParam
      ? getNextPageParam(previousPageData, swr.data || [])
      : previousPageData?.nextCursor

    if (!nextParam) return null

    const query = new URLSearchParams({
      ...initialParams,
      cursor: String(nextParam),
      limit: String(pageSize),
    }).toString()

    return `${baseUrl}?${query}`
  }

  const swr = useSWRInfinite<T>(getKey, fetcher, {
    revalidateFirstPage: false,
    ...swrOptions,
  })

  const list = swr.data
    ? swr.data.flatMap((item: any) =>
        Array.isArray(item?.data) ? item.data : item,
      )
    : []

  const loading = !swr.data && !swr.error

  const loadingMore =
    loading ||
    (swr.size > 0 && swr.data && typeof swr.data[swr.size - 1] === 'undefined')

  const isEmpty = swr.data?.[0]?.data?.length === 0
  const lastPage = swr.data?.[swr.data.length - 1]

  const noMore = isEmpty || (isEnd ? isEnd(lastPage) : !lastPage?.nextCursor)

  return {
    ...swr,
    list,
    loading,
    loadingMore,
    noMore,
    loadMore: () => swr.setSize(swr.size + 1),
  }
}

// 简易封装
type InfiniteResponse<T> = {
  data: T[]
  nextCursor?: string
}

export function useGetRequestInfinite3<T>(url: string, pageSize = 10) {
  const getKey = (
    pageIndex: number,
    previousPageData: InfiniteResponse<T> | null,
  ) => {
    // 🚫 没有下一页了
    if (previousPageData && !previousPageData.nextCursor) return null
    // ✅ 第一页
    if (pageIndex === 0) {
      return `${url}?${[LIMIT_NAME]}=${pageSize}`
    }
    // ✅ 后续页（真正 cursor）
    const query = new URLSearchParams({
      [CURSOR_NAME]: String(pageIndex),
      [LIMIT_NAME]: String(LIMIT_SIZE),
    }).toString()
    return `${url}?${query}`
  }

  const swr = useSWRInfinite<InfiniteResponse<T>>(getKey, fetcher)

  return {
    ...swr,

    // 🔥 扁平化数据（非常重要）
    list: swr.data ? swr.data.flatMap((page) => page.data) : [],

    // 🔥 是否还有更多
    hasMore: swr.data ? !!swr.data[swr.data.length - 1]?.nextCursor : true,
  }
}

export function useGetRequestInfinite4<T>(url: string) {
  const getKey = (pageIndex: number, previousPageData: any) => {
    // reached the end
    if (previousPageData && !previousPageData.data) return null
    // first page, we don't have `previousPageData`

    // 第一页
    if (pageIndex === 0) {
      const query = new URLSearchParams({
        limit: String(LIMIT_SIZE),
      }).toString()
      return `${url}?${query}`
    }

    const query = new URLSearchParams({
      cursor: String(pageIndex),
      limit: String(LIMIT_SIZE),
    }).toString()
    // add the cursor to the API endpoint
    return `${url}?${query}`
  }
  return useSWRInfinite<T>(getKey, fetcher)
}
