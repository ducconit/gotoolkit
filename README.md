# 🚀 GoToolkit

[![Go Reference](https://img.shields.io/badge/go-1.26+-00ADD8.svg?style=flat-square&logo=go)](https://golang.org)
[![GitHub Profile](https://img.shields.io/badge/author-DNT-orange.svg?style=flat-square)](https://github.com/ducconit)

Bộ công cụ (toolkit) phát triển ứng dụng bằng ngôn ngữ **Golang 1.26+**, được xây dựng với triết lý tối thượng: **Cực kỳ tối giản (Zero-boilerplate), hiệu năng thô tối đa (Zero-allocation ở hot-path) và trải nghiệm lập trình viên (DX) tốt nhất.**

Bộ toolkit tận dụng tối đa các tính năng hiện đại của Go (Generics nâng cao, Iterators chuẩn Go 1.23+, Structured Logging với `log/slog`,...) để giúp bạn xây dựng dịch vụ web nhanh hơn và nhàn hơn.

---

## 📦 Installation

Để cài đặt bộ toolkit vào dự án Go của bạn, chạy lệnh sau:

```bash
go get github.com/ducconit/gotoolkit
```

---

## 📂 Danh sách các Package & Tính năng

Bộ toolkit được cấu trúc thành các module độc lập đặt ngay tại thư mục root, cho phép bạn dễ dàng import chỉ những gì cần thiết:

### 1. 🔒 [secureapi](./secureapi) - Application-Level Encryption
Giải pháp bảo mật API mức ứng dụng (ALE) giúp mã hóa toàn bộ dữ liệu trao đổi giữa Client (SPA/Mobile) và Server để chống F12, chống xem lén dữ liệu trên đường truyền (MitM).
*   Trao đổi khóa động bảo mật qua **ECDH (Curve P-256)**.
*   Mã hóa đối xứng hiệu năng cao qua **AES-256-GCM**.
*   Tự động hóa bằng Middleware cho cả **stdlib `net/http`** và **Gin framework**.
*   Bảo vệ bộ nhớ nâng cao (sync.Pool giới hạn buffer, chống rò rỉ session keys).

### 2. 🔑 [encrypt](./encrypt) - Crypto Helpers
Các tiện ích mã hóa đối xứng an toàn dùng chung cho toàn bộ dự án (mã hóa cookie, database, session...).
*   Cung cấp `encrypt.AESGCM(...)` để mã hóa và `encrypt.DecryptAESGCM(...)` giải mã an toàn.
*   Hỗ trợ `encrypt.DecryptAESGCMInPlace(...)` giải mã trực tiếp trên slice (Zero-allocation) cho các tác vụ cần tối ưu hóa bộ nhớ khắt khe.

### 👥 3. [rbac](./rbac) - Role-Based Access Control
Quản lý và kiểm tra phân quyền người dùng theo vai trò (Role) và quyền hạn (Permission) một cách trực quan, mạch lạc.

### 📝 4. [str](./str) - String Utilities
Các helper xử lý chuỗi ký tự tối ưu hóa hiệu năng, đặc biệt hỗ trợ tốt cho tiếng Việt:
*   `RemoveAccents`: Loại bỏ dấu tiếng Việt chuẩn xác (hỗ trợ cả Unicode dựng sẵn và tổ hợp).
*   `Slugify`: Tạo slug thân thiện với SEO cho URL.
*   `CleanSpace`: Dọn dẹp khoảng trắng thừa (spaces, tabs, newlines).
*   `ContainsHTML`: Phát hiện mã độc HTML/JS để chống tấn công XSS.

### 🎲 5. [random](./random) - Secure Random Generator
Tạo các giá trị ngẫu nhiên có độ an toàn mã hóa cao (cryptographically secure), tránh sử dụng `math/rand` không an toàn:
*   Sinh chuỗi ngẫu nhiên (string tokens), số ngẫu nhiên.
*   Sinh UUID v4.

### 🎛️ 6. [feature](./feature) - Feature Flags
Hệ thống bật/tắt tính năng động (Feature toggles) giúp deploy các tính năng mới an toàn mà không cần khởi động lại ứng dụng.

### ⚙️ 7. [sys](./sys) - System Utilities
Các tiện ích hệ thống giúp kiểm tra môi trường, giám sát tài nguyên CPU/RAM, xử lý graceful shutdown cho ứng dụng.

---

## 🛠️ Quy tắc Thiết kế & Phát triển

Nếu bạn muốn đóng góp code hoặc phát triển thêm các package mới cho `gotoolkit`, vui lòng tuân thủ các nguyên tắc thiết kế sau:
1.  **Sử dụng Go hiện đại**: Viết code tương thích với **Go 1.26+**, tận dụng `iter.Seq`, `slices`, `maps`, `log/slog` thay vì các thư viện bên ngoài.
2.  **Zero-Boilerplate**: Code của thư viện phải tối giản tối đa để người dùng chỉ cần viết vài dòng là chạy được.
3.  **Hot Path Optimization**: Hướng tới **0 allocs/op** ở các tác vụ đọc (read-heavy) bằng cách tái sử dụng bộ nhớ (`sync.Pool`) và tối ưu hóa Escape Analysis.
4.  **Table-Driven Tests**: Mọi tính năng mới bắt buộc phải có unit test đầy đủ sử dụng cấu trúc Table-Driven Tests của Go.

---

## 🤝 Đóng góp (Contributing)

Mọi ý kiến đóng góp, báo lỗi hoặc yêu cầu tính năng mới đều được chào đón! Bạn có thể:
1. Fork dự án này.
2. Tạo nhánh tính năng mới (`git checkout -b feature/amazing-feature`).
3. Commit các thay đổi của bạn (`git commit -m 'Add some amazing feature'`).
4. Push lên nhánh vừa tạo (`git push origin feature/amazing-feature`).
5. Mở một Pull Request.

---

## 📈 Star History

[![Star History Chart](https://api.star-history.com/svg?repos=ducconit/gotoolkit&type=Date)](https://star-history.com/#ducconit/gotoolkit&Date)

---

## 👤 Owner & Contact

*   **Owner**: DNT <ducconit@gmail.com>
*   **GitHub**: [github.com/ducconit](https://github.com/ducconit)
