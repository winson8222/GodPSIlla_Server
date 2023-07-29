/*
  Warnings:

  - Added the required column `idlname` to the `Version` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "Version" ADD COLUMN     "idlname" TEXT NOT NULL;
