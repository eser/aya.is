# AYA Web Sitesi

[![Discord](https://img.shields.io/discord/1072074800622739476?color=7289da&logo=discord&logoColor=white)](https://discord.gg/itdepremyardim)
[![GitHub issues](https://img.shields.io/github/issues/eser/aya.is)](https://github.com/eser/aya.is/issues)

(For English please click [here](README.en.md))

Açık kaynak yazılımı ve açık kaynak verilerine ortak ilgi duyan bireylerden oluşan sıkı bağlı ve tutkulu bir topluluğuz.
Toplumun iyileştirilmesi için teknolojiden yararlanma kararlılığımızla birleşerek, Türkiye'de meydana gelen yıkıcı
[Şubat 2023 depremi](https://en.wikipedia.org/wiki/2023_Turkey%E2%80%93Syria_earthquake) sonrasında güçlerimizi
birleştirdik.

Depremin ardından, etkilenenlere destek olma ve aileleri ile sevdiklerini yeniden bir araya getirme sürecine yardım etme
acil ihtiyacını fark ettik. Bu amaç doğrultusunda, bireyler arasındaki arama ve bağlantıyı kolaylaştırmayı hedefleyen
yenilikçi araçlar ve çözümler geliştirmek için yola çıktık.

Yorulmak bilmeyen iş birliği ve kolektif çabalar sayesinde topluluğumuz ilham verici bir büyüklüğe ulaştı ve 24.000
kişiden oluşan bir topluluk haline geldik. Birlikte, aileleri yeniden bir araya getirmeye, umudu geri kazandırmaya ve
depremden etkilenenlere teselli sunmaya yardımcı olmak için teknolojinin gücünden yararlanarak araçlarımızı sürekli
olarak geliştiriyoruz ve genişletiyoruz.

İlerlemeye devam ederken, taahhüdümüzden sapmadan devam ediyoruz. Araçlarımızı geliştirmeye, açık kaynak projelere aktif
olarak katkıda bulunmaya ve bilgi paylaşımı ve iş birliğinin geliştiği destekleyici bir ortamı teşvik etmeye
çalışıyoruz. Açık kaynak topluluğunun bir parçası olmaktan gurur duyuyoruz ve sizlerle olan yolculuğumuza dört gözle
devam etmek için sabırsızlanıyoruz.

## Görevimiz

Açık kaynak çözümlerini kullanarak, bilgi sistemlerini uygulayarak ve mühendislik uygulamalarını hayata geçirerek,
içinde yaşadığımız toplumun iyileştirilmesine aktif olarak katkıda bulunmayı hedefliyoruz. Öncelikli odak noktamız,
sosyal sorumluluk ihtiyaçlarını karşılamak ve bu yollarla çeşitli toplumsal zorluklarla baş etmektir. Uzmanlığımızı ve
kaynaklarımızı kullanarak, topluluğumuz içinde pozitif bir etki yaratmayı ve sürdürülebilir kalkınmayı teşvik etmeyi
amaçlıyoruz.

## Teknoloji

Bu bir Docker Compose tabanlı monorepo projesidir. Aşağıdaki bileşenleri içerir:

- **Frontend (webclient)**: Next.js ve Shadcn UI kullanır
- **Backend (services)**: Go dilinde yazılmış REST API servisleri
- **Veritabanı**: PostgreSQL

Önceden kurulu olması gerekenler:

- [Docker](https://docker.com) (Orbstack önerilir)
- [Make](https://www.gnu.org/software/make/) (Unix/macOS sistemlerde genellikle yüklüdür)
- [Git](https://git-scm.com/)

## Kurulum ve Başlangıç

GitHub repository'sini klonlayın:

```bash
$ git clone git@github.com:eser/aya.is.git
$ cd aya.is
```

Backend servisi için gerekli ortam değişkenlerini ayarlayın:

```bash
# .env.local dosyası oluşturun (bu dosya git'e eklenmez)
$ cat > apps/services/.env.local << EOF
AUTH__JWT_SECRET=herhangi-bir-gizli-anahtar
EOF
```

> **Not:** `AUTH__JWT_SECRET` kimlik doğrulama için zorunludur. Geliştirme ortamında herhangi bir değer kullanabilirsiniz.

Projeyi başlatmak için Make komutlarını kullanın:

```bash
# Tüm servisleri build et ve başlat
$ make up
```

Diğer kullanışlı Make komutları:

```bash
$ make help      # Tüm komutları göster
$ make logs      # Container loglarını göster
$ make stop      # Servisleri durdur
$ make restart   # Servisleri yeniden başlat
$ make down      # Containerları tamamen kaldır
```

## Proje Yönetimi ve CLI

Servis container'ına bağlanmak için:

```bash
$ make cli
```

Bu komut sizi backend servisinin bash shell'ine bağlar. Buradan veritabanı yönetimi ve diğer işlemleri yapabilirsiniz.

### Örnek

Profil getirme:

```js
await backend.getProfile("tr", "eser");
```

### Veri Modeli

Proje aşağıdaki ana veri yapılarını kullanır:

- **Profile**: Kullanıcı profilleri (bireyler, organizasyonlar, topluluklar)
- **Story**: İçerik ve yazılar (blog yazıları, haberler, etkinlikler)
- **User**: Sistem kullanıcıları ve kimlik doğrulama bilgileri

## Nasıl Katkıda Bulunabilirsin?

Herkesten katkı bekliyoruz. Başlamak için lütfen [katkıda bulunma kılavuzumuzu](CONTRIBUTING.md) okuyun. Yardım etmek
isterseniz [issue'lara](https://github.com/eser/aya.is/issues) göz atabilirsiniz. Herhangi bir sorunuz varsa
[Discord sunucumuza](https://discord.gg/itdepremyardim) katılmaktan çekinmeyin. Herhangi bir noktada takılırsanız,
GitHub Issues veya Discord üzerinden yardım istemekten çekinmeyin.

## Lisans

Apache 2.0, detaylar için [LICENSE](LICENSE) dosyasını inceleyebilirsiniz.
