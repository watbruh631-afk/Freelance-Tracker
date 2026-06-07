package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

const MAKS_PROYEK = 100

type StatusProyek int

const (
	StatusBerjalan StatusProyek = iota
	StatusSelesai
	StatusPending
)

type Proyek struct {
	ID           int
	NamaProyek   string
	NamaKlien    string
	Deskripsi    string
	Status       StatusProyek
	Bayaran      float64
	Deadline     string
	TanggalMulai string
	Aktif        bool
}

type DaftarProyek [MAKS_PROYEK]Proyek

var daftarProyek DaftarProyek
var jumlahProyek int
var counterID int

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func bacaString(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func bacaInt(prompt string) int {
	var nilai int
	fmt.Print(prompt)
	fmt.Scan(&nilai)
	bufio.NewReader(os.Stdin).ReadString('\n')
	return nilai
}

func bacaFloat(prompt string) float64 {
	var nilai float64
	fmt.Print(prompt)
	fmt.Scan(&nilai)
	bufio.NewReader(os.Stdin).ReadString('\n')
	return nilai
}

func tekanEnter() {
	fmt.Print("\n  Tekan Enter untuk kembali...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func cetakGaris() {
	fmt.Println(strings.Repeat("=", 56))
}

func cetakHeader(judul string) {
	cetakGaris()
	pad := (56 - len(judul)) / 2
	fmt.Println(strings.Repeat(" ", pad) + judul)
	cetakGaris()
}

func cetakHeaderTabel() {
	fmt.Printf("  %-4s | %-24s | %-18s | %-18s | %-14s | %-10s\n",
		"No", "Nama Proyek", "Klien", "Status", "Bayaran (Rp)", "Deadline")
	fmt.Println("  " + strings.Repeat("-", 105))
}

func tampilkanBarisTabel(p Proyek, nomor int) {
	fmt.Printf("  %-4d | %-24s | %-18s | %-18s | %-14.0f | %-10s\n",
		nomor, p.NamaProyek, p.NamaKlien, statusKeString(p.Status), p.Bayaran, p.Deadline)
}

func tampilkanDetailProyek(p Proyek) {
	cetakGaris()
	fmt.Printf("  ID            : %d\n  Nama Proyek   : %s\n  Klien         : %s\n",
		p.ID, p.NamaProyek, p.NamaKlien)
	fmt.Printf("  Deskripsi     : %s\n  Status        : %s\n  Bayaran       : Rp %.0f\n",
		p.Deskripsi, statusKeString(p.Status), p.Bayaran)
	fmt.Printf("  Tgl Mulai     : %s\n  Deadline      : %s\n", p.TanggalMulai, p.Deadline)
	cetakGaris()
}

func statusKeString(s StatusProyek) string {
	if s == StatusBerjalan {
		return "Sedang Dikerjakan"
	}
	if s == StatusSelesai {
		return "Selesai"
	}
	return "Pending"
}

func intKeStatus(p int) StatusProyek {
	if p == 1 {
		return StatusBerjalan
	}
	if p == 2 {
		return StatusSelesai
	}
	return StatusPending
}

func validasiTanggal(t string) bool {
	return len(t) == 10 && t[2] == '/' && t[5] == '/'
}

func tanggalHariIni() string {
	t := time.Now()
	return fmt.Sprintf("%02d/%02d/%04d", t.Day(), int(t.Month()), t.Year())
}

func bangunIdxAktif(idx *[MAKS_PROYEK]int) int {
	n := 0
	i := 0
	for i < jumlahProyek {
		if daftarProyek[i].Aktif {
			idx[n] = i
			n++
		}
		i++
	}
	return n
}

func tukar(arr *[MAKS_PROYEK]int, i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

func sequentialSearch(keyword, field string, hasil *[MAKS_PROYEK]int, jml *int) {
	*jml = 0
	kw := strings.ToLower(keyword)
	i := 0
	for i < jumlahProyek {
		if daftarProyek[i].Aktif {
			val := ""
			if field == "proyek" {
				val = daftarProyek[i].NamaProyek
			} else {
				val = daftarProyek[i].NamaKlien
			}
			if strings.Contains(strings.ToLower(val), kw) {
				hasil[*jml] = i
				*jml++
			}
		}
		i++
	}
}

func binarySearchID(id int, idx *[MAKS_PROYEK]int, n int) int {
	i := 1
	for i < n {
		kunci := idx[i]
		j := i - 1
		for j >= 0 && daftarProyek[idx[j]].ID > daftarProyek[kunci].ID {
			idx[j+1] = idx[j]
			j--
		}
		idx[j+1] = kunci
		i++
	}
	kiri, kanan, hasil := 0, n-1, -1
	for kiri <= kanan && hasil == -1 {
		tengah := (kiri + kanan) / 2
		if daftarProyek[idx[tengah]].ID == id {
			hasil = idx[tengah]
		} else if daftarProyek[idx[tengah]].ID < id {
			kiri = tengah + 1
		} else {
			kanan = tengah - 1
		}
	}
	return hasil
}

func selectionSort(idx *[MAKS_PROYEK]int, n int, field string, asc bool) {
	i := 0
	for i < n-1 {
		ext := i
		j := i + 1
		for j < n {
			var a, b float64
			if field == "bayaran" {
				a = daftarProyek[idx[j]].Bayaran
				b = daftarProyek[idx[ext]].Bayaran
			} else {
				if daftarProyek[idx[j]].Deadline < daftarProyek[idx[ext]].Deadline {
					a, b = 0, 1
				} else {
					a, b = 1, 0
				}
			}
			if (asc && a < b) || (!asc && a > b) {
				ext = j
			}
			j++
		}
		if ext != i {
			tukar(idx, i, ext)
		}
		i++
	}
}

func insertionSort(idx *[MAKS_PROYEK]int, n int, field string, asc bool) {
	i := 1
	for i < n {
		kunci := idx[i]
		j := i - 1
		selesai := false
		for j >= 0 && !selesai {
			var a, b string
			if field == "proyek" {
				a = strings.ToLower(daftarProyek[idx[j]].NamaProyek)
				b = strings.ToLower(daftarProyek[kunci].NamaProyek)
			} else {
				a = strings.ToLower(daftarProyek[idx[j]].NamaKlien)
				b = strings.ToLower(daftarProyek[kunci].NamaKlien)
			}
			if (asc && a > b) || (!asc && a < b) {
				idx[j+1] = idx[j]
				j--
			} else {
				selesai = true
			}
		}
		idx[j+1] = kunci
		i++
	}
}

func tampilTabel(idx [MAKS_PROYEK]int, n int) {
	cetakHeaderTabel()
	i := 0
	for i < n {
		tampilkanBarisTabel(daftarProyek[idx[i]], i+1)
		i++
	}
	fmt.Println("  " + strings.Repeat("-", 105))
}

func inputProyek() Proyek {
	var p Proyek
	counterID++
	p.ID = counterID
	p.Aktif = true
	p.TanggalMulai = tanggalHariIni()
	p.NamaProyek = bacaString("  Nama Proyek          : ")
	p.NamaKlien = bacaString("  Nama Klien           : ")
	p.Deskripsi = bacaString("  Deskripsi Proyek     : ")
	fmt.Println("  Status: 1.Berjalan  2.Selesai  3.Pending")
	p.Status = intKeStatus(bacaInt("  Pilihan Status       : "))
	p.Bayaran = bacaFloat("  Bayaran (Rp)          : ")
	fmt.Printf("  Tgl Mulai             : %s (otomatis)\n", p.TanggalMulai)
	valid := false
	for !valid {
		p.Deadline = bacaString("  Deadline (DD/MM/YYYY) : ")
		valid = validasiTanggal(p.Deadline)
		if !valid {
			fmt.Println("  [!] Format tidak valid. Contoh: 31/12/2026")
		}
	}
	return p
}

func tambahProyek() {
	clearScreen()
	cetakHeader("TAMBAH PROYEK BARU")
	fmt.Println()
	if jumlahProyek >= MAKS_PROYEK {
		fmt.Println("  [!] Kapasitas penuh!")
		tekanEnter()
		return
	}
	p := inputProyek()
	daftarProyek[jumlahProyek] = p
	jumlahProyek++
	fmt.Println()
	cetakGaris()
	fmt.Printf("  [OK] Proyek \"%s\" ditambahkan! ID: %d\n", p.NamaProyek, p.ID)
	cetakGaris()
	tekanEnter()
}

func lihatSemuaProyek() {
	clearScreen()
	cetakHeader("DAFTAR SEMUA PROYEK")
	var idx [MAKS_PROYEK]int
	n := bangunIdxAktif(&idx)
	if n == 0 {
		fmt.Println("\n  Belum ada proyek.")
		tekanEnter()
		return
	}
	fmt.Println("\n  Urutkan: 1.Deadline[SS] 2.Bayaran[SS] 3.NamaProyek[IS] 4.NamaKlien[IS] 5.Tidak")
	pilUrut := bacaInt("  Pilihan: ")
	asc := true
	if pilUrut >= 1 && pilUrut <= 4 {
		pilArah := bacaInt("  Arah (1=Asc 2=Desc): ")
		asc = (pilArah != 2)
	}
	if pilUrut == 1 {
		selectionSort(&idx, n, "deadline", asc)
	} else if pilUrut == 2 {
		selectionSort(&idx, n, "bayaran", asc)
	} else if pilUrut == 3 {
		insertionSort(&idx, n, "proyek", asc)
	} else if pilUrut == 4 {
		insertionSort(&idx, n, "klien", asc)
	}
	fmt.Println()
	tampilTabel(idx, n)
	fmt.Printf("\n  Total: %d proyek\n", n)
	tekanEnter()
}

func cariProyek() {
	clearScreen()
	cetakHeader("CARI PROYEK")
	fmt.Println("\n  1.Nama Proyek [Sequential]  2.Nama Klien [Sequential]  3.ID [Binary]")
	pilMetode := bacaInt("  Pilihan: ")
	fmt.Println()

	var hasilIdx [MAKS_PROYEK]int
	var jmlHasil int

	if pilMetode == 1 || pilMetode == 2 {
		field := "proyek"
		prompt := "  Kata kunci nama proyek : "
		if pilMetode == 2 {
			field = "klien"
			prompt = "  Kata kunci nama klien  : "
		}
		kw := bacaString(prompt)
		sequentialSearch(kw, field, &hasilIdx, &jmlHasil)
		fmt.Printf("\n  [Sequential Search] keyword: \"%s\"\n\n", kw)
		if jmlHasil == 0 {
			fmt.Println("  Tidak ada proyek yang cocok.")
		} else {
			tampilTabel(hasilIdx, jmlHasil)
			fmt.Printf("  Ditemukan: %d proyek\n", jmlHasil)
		}
	} else if pilMetode == 3 {
		idCari := bacaInt("  ID Proyek yang dicari : ")
		var aktifIdx [MAKS_PROYEK]int
		n := bangunIdxAktif(&aktifIdx)
		idxHasil := binarySearchID(idCari, &aktifIdx, n)
		fmt.Printf("\n  [Binary Search] ID: %d\n\n", idCari)
		if idxHasil == -1 {
			fmt.Printf("  ID %d tidak ditemukan.\n", idCari)
		} else {
			tampilkanDetailProyek(daftarProyek[idxHasil])
		}
	} else {
		fmt.Println("  [!] Pilihan tidak valid.")
	}
	tekanEnter()
}

func cariBinaryUntukCRUD(idInput int) int {
	var aktifIdx [MAKS_PROYEK]int
	n := bangunIdxAktif(&aktifIdx)
	return binarySearchID(idInput, &aktifIdx, n)
}

func editProyek() {
	clearScreen()
	cetakHeader("EDIT PROYEK")
	idEdit := bacaInt("\n  ID Proyek yang akan diedit: ")
	idxHasil := cariBinaryUntukCRUD(idEdit)
	if idxHasil == -1 {
		fmt.Printf("\n  [Binary Search] ID %d tidak ditemukan.\n", idEdit)
		tekanEnter()
		return
	}
	fmt.Println("\n  [Binary Search] Ditemukan!")
	tampilkanDetailProyek(daftarProyek[idxHasil])
	fmt.Println("  (Kosongkan untuk tidak mengubah)")

	if v := bacaString("  Nama Proyek Baru : "); v != "" {
		daftarProyek[idxHasil].NamaProyek = v
	}
	if v := bacaString("  Nama Klien Baru  : "); v != "" {
		daftarProyek[idxHasil].NamaKlien = v
	}
	if v := bacaString("  Deskripsi Baru   : "); v != "" {
		daftarProyek[idxHasil].Deskripsi = v
	}
	fmt.Println("  Status (0=tidak ubah): 1.Berjalan 2.Selesai 3.Pending")
	if p := bacaInt("  Pilihan: "); p >= 1 && p <= 3 {
		daftarProyek[idxHasil].Status = intKeStatus(p)
	}
	if v := bacaFloat("  Bayaran Baru (0=tidak ubah): Rp "); v > 0 {
		daftarProyek[idxHasil].Bayaran = v
	}
	if v := bacaString("  Deadline Baru (kosong=tidak ubah): "); v != "" {
		if validasiTanggal(v) {
			daftarProyek[idxHasil].Deadline = v
		} else {
			fmt.Println("  [!] Format tidak valid, deadline tidak diubah.")
		}
	}
	cetakGaris()
	fmt.Println("  [OK] Proyek berhasil diperbarui!")
	cetakGaris()
	tekanEnter()
}

func hapusProyek() {
	clearScreen()
	cetakHeader("HAPUS PROYEK")
	idHapus := bacaInt("\n  ID Proyek yang akan dihapus: ")
	idxHasil := cariBinaryUntukCRUD(idHapus)
	if idxHasil == -1 {
		fmt.Printf("\n  [Binary Search] ID %d tidak ditemukan.\n", idHapus)
		tekanEnter()
		return
	}
	fmt.Println("\n  [Binary Search] Ditemukan!")
	tampilkanDetailProyek(daftarProyek[idxHasil])
	if strings.ToLower(bacaString("  Konfirmasi hapus? (y/n): ")) == "y" {
		nama := daftarProyek[idxHasil].NamaProyek
		daftarProyek[idxHasil].Aktif = false
		cetakGaris()
		fmt.Printf("  [OK] Proyek \"%s\" dihapus!\n", nama)
		cetakGaris()
	} else {
		fmt.Println("\n  Dibatalkan.")
	}
	tekanEnter()
}

func updateStatusProyek() {
	clearScreen()
	cetakHeader("UPDATE STATUS PROYEK")
	idUpdate := bacaInt("\n  ID Proyek: ")
	idxHasil := -1
	i := 0
	for i < jumlahProyek {
		if daftarProyek[i].Aktif && daftarProyek[i].ID == idUpdate {
			idxHasil = i
			i = jumlahProyek
		} else {
			i++
		}
	}
	if idxHasil == -1 {
		fmt.Printf("\n  [Sequential Search] ID %d tidak ditemukan.\n", idUpdate)
		tekanEnter()
		return
	}
	fmt.Printf("\n  [Sequential Search] Ditemukan!\n  Nama  : %s\n  Status: %s\n\n",
		daftarProyek[idxHasil].NamaProyek, statusKeString(daftarProyek[idxHasil].Status))
	fmt.Println("  1.Sedang Dikerjakan  2.Selesai  3.Pending")
	if p := bacaInt("\n  Pilihan: "); p >= 1 && p <= 3 {
		lama := statusKeString(daftarProyek[idxHasil].Status)
		daftarProyek[idxHasil].Status = intKeStatus(p)
		cetakGaris()
		fmt.Printf("  [OK] %s --> %s\n", lama, statusKeString(daftarProyek[idxHasil].Status))
		cetakGaris()
	} else {
		fmt.Println("  [!] Pilihan tidak valid.")
	}
	tekanEnter()
}

func lihatLaporan() {
	clearScreen()
	cetakHeader("LAPORAN PROYEK FREELANCE")
	var idx [MAKS_PROYEK]int
	n := bangunIdxAktif(&idx)
	if n == 0 {
		fmt.Println("\n  Belum ada data.")
		tekanEnter()
		return
	}

	jmlB, jmlS, jmlP := 0, 0, 0
	totalS, totalA := 0.0, 0.0
	i := 0
	for i < jumlahProyek {
		if daftarProyek[i].Aktif {
			if daftarProyek[i].Status == StatusBerjalan {
				jmlB++
			} else if daftarProyek[i].Status == StatusSelesai {
				jmlS++
				totalS += daftarProyek[i].Bayaran
			} else {
				jmlP++
			}
			totalA += daftarProyek[i].Bayaran
		}
		i++
	}
	fmt.Println()
	cetakGaris()
	fmt.Printf("  Total: %d | Berjalan: %d | Selesai: %d | Pending: %d\n", n, jmlB, jmlS, jmlP)
	fmt.Printf("  Pendapatan Selesai : Rp %.0f\n  Potensi Total      : Rp %.0f\n", totalS, totalA)
	cetakGaris()

	fmt.Println("\n  1.Berjalan 2.Selesai 3.Pending 4.Deadline[SS] 5.Bayaran[SS] 6.NamaAZ[IS] 7.NamaZA[IS] 0.Kembali")
	pil := bacaInt("  Pilihan: ")
	if pil == 0 {
		return
	}
	fmt.Println()

	if pil >= 1 && pil <= 3 {
		statusMap := [4]StatusProyek{0, StatusBerjalan, StatusSelesai, StatusPending}
		judulMap := [4]string{"", "PROYEK SEDANG DIKERJAKAN", "PROYEK SELESAI", "PROYEK PENDING"}
		cetakHeader(judulMap[pil])
		fmt.Println()
		var fIdx [MAKS_PROYEK]int
		jf := 0
		i := 0
		for i < n {
			if daftarProyek[idx[i]].Status == statusMap[pil] {
				fIdx[jf] = idx[i]
				jf++
			}
			i++
		}
		if jf == 0 {
			fmt.Println("  Tidak ada proyek dengan status ini.")
		} else {
			tampilTabel(fIdx, jf)
			total := 0.0
			k := 0
			for k < jf {
				total += daftarProyek[fIdx[k]].Bayaran
				k++
			}
			fmt.Printf("  Jumlah: %d | Total Bayaran: Rp %.0f\n", jf, total)
		}
	} else if pil == 4 {
		cetakHeader("DEADLINE TERDEKAT [Selection Sort]")
		fmt.Println()
		selectionSort(&idx, n, "deadline", true)
		tampilTabel(idx, n)
	} else if pil == 5 {
		cetakHeader("BAYARAN TERTINGGI [Selection Sort]")
		fmt.Println()
		selectionSort(&idx, n, "bayaran", false)
		tampilTabel(idx, n)
	} else if pil == 6 {
		cetakHeader("NAMA A-Z [Insertion Sort]")
		fmt.Println()
		insertionSort(&idx, n, "proyek", true)
		tampilTabel(idx, n)
	} else if pil == 7 {
		cetakHeader("NAMA Z-A [Insertion Sort]")
		fmt.Println()
		insertionSort(&idx, n, "proyek", false)
		tampilTabel(idx, n)
	} else {
		fmt.Println("  [!] Pilihan tidak valid.")
	}
	tekanEnter()
}

func isiDataContoh() {
	contoh := [1]Proyek{
		{ID: 1, NamaProyek: "Website E-Commerce", NamaKlien: "PT Maju Jaya",
			Deskripsi: "Toko online + sistem pembayaran", Status: StatusBerjalan,
			Bayaran: 5000000, TanggalMulai: "01/05/2026", Deadline: "30/06/2026", Aktif: true},
	}
	i := 0
	for i < 1 {
		daftarProyek[i] = contoh[i]
		i++
	}
	jumlahProyek = 1
	counterID = 1
}

func tampilkanMenuUtama() {
	clearScreen()
	jmlAktif := 0
	i := 0
	for i < jumlahProyek {
		if daftarProyek[i].Aktif {
			jmlAktif++
		}
		i++
	}
	cetakGaris()
	fmt.Println("     APLIKASI MANAJEMEN TRACKING FREELANCE")
	fmt.Println("             Algoritma & Pemrograman 2")
	cetakGaris()
	fmt.Printf("  Tanggal : %s  |  Proyek Aktif : %d\n", tanggalHariIni(), jmlAktif)
	cetakGaris()
	fmt.Println()
	fmt.Println("  1. Lihat Semua Proyek    5. Update Status")
	fmt.Println("  2. Tambah Proyek         6. Cari Proyek")
	fmt.Println("  3. Edit Proyek           7. Laporan & Statistik")
	fmt.Println("  4. Hapus Proyek          0. Keluar")
	fmt.Println()
	cetakGaris()
}

func prosesMenu(pilihan int) bool {
	if pilihan == 0 {
		return true
	}
	if pilihan == 1 {
		lihatSemuaProyek()
	} else if pilihan == 2 {
		tambahProyek()
	} else if pilihan == 3 {
		editProyek()
	} else if pilihan == 4 {
		hapusProyek()
	} else if pilihan == 5 {
		updateStatusProyek()
	} else if pilihan == 6 {
		cariProyek()
	} else if pilihan == 7 {
		lihatLaporan()
	} else {
		clearScreen()
		fmt.Println("  [!] Pilihan tidak valid.")
		tekanEnter()
	}
	return false
}

func main() {
	jumlahProyek = 0
	counterID = 0
	isiDataContoh()
	keluar := false
	for !keluar {
		tampilkanMenuUtama()
		keluar = prosesMenu(bacaInt("  Pilihan: "))
	}
	clearScreen()
	cetakGaris()
	fmt.Println("  Terima kasih! Aplikasi Manajemen Freelance.")
	cetakGaris()
	fmt.Println()
}
