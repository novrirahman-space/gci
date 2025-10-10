# Debugging Explanation

## Identifikasi Masalah Utama
Kode awal tidak dapat diandalkan karena:
1. Tidak ada sinkronisasi antar goroutine dengan main thread.
2. Mengandalkan `time.Sleep()` untuk menunggu goroutine.
3. Potensi goroutine leak dan koneksi HTTP tidak ditutup.
4. Tidak ada batas waktu untuk operasi jaringan.

Akibatnya hasil acak dan sering hilang sebelum semua proses selesai.

## Solusi Idiomatis di Go
Gunakan `sync.WaitGroup` untuk menunggu semua goroutine selesai secara deterministik.

Pola umum:
```go
var wg sync.WaitGroup
wg.Add(n)
for range n {
    go func() {
        defer wg.Done()
        // pekerjaan
    }()
}
wg.Wait()
