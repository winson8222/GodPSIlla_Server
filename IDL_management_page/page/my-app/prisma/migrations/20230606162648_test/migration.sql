-- CreateTable
CREATE TABLE "Microservice" (
    "id" TEXT NOT NULL,
    "svcname" TEXT NOT NULL,

    CONSTRAINT "Microservice_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Version" (
    "id" TEXT NOT NULL,
    "vname" TEXT NOT NULL,
    "upstreamurl" TEXT NOT NULL,
    "idlfile" BYTEA NOT NULL,
    "microserviceId" TEXT,

    CONSTRAINT "Version_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "Microservice_svcname_key" ON "Microservice"("svcname");

-- AddForeignKey
ALTER TABLE "Version" ADD CONSTRAINT "Version_microserviceId_fkey" FOREIGN KEY ("microserviceId") REFERENCES "Microservice"("id") ON DELETE SET NULL ON UPDATE CASCADE;
