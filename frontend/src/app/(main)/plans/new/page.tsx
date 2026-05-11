"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import {
  Box,
  Button,
  Card,
  Heading,
  Input,
  Text,
  VStack,
} from "@chakra-ui/react";

export default function NewPlanPage() {
  const router = useRouter();
  const [title, setTitle] = useState("");
  const [totalPages, setTotalPages] = useState("");
  const [targetDate, setTargetDate] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // TODO: API接続（フェーズ3で実装）
    router.push("/plans");
  };

  return (
    <VStack gap={6} align="stretch" maxW="600px">
      <Heading size="lg">新しい計画を作成</Heading>

      <Card.Root>
        <Card.Body>
          <Box as="form" onSubmit={handleSubmit}>
            <VStack gap={4}>
              <VStack gap={1} align="start" w="full">
                <Text fontSize="sm" fontWeight="medium">教材名</Text>
                <Input
                  placeholder="例: AWS SAA対策本"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  required
                  maxLength={200}
                />
              </VStack>

              <VStack gap={1} align="start" w="full">
                <Text fontSize="sm" fontWeight="medium">総ページ数</Text>
                <Input
                  type="number"
                  placeholder="例: 300"
                  value={totalPages}
                  onChange={(e) => setTotalPages(e.target.value)}
                  required
                  min={1}
                  max={10000}
                />
              </VStack>

              <VStack gap={1} align="start" w="full">
                <Text fontSize="sm" fontWeight="medium">目標期限</Text>
                <Input
                  type="date"
                  value={targetDate}
                  onChange={(e) => setTargetDate(e.target.value)}
                  required
                />
              </VStack>

              <Button type="submit" colorPalette="blue" w="full">
                計画を作成
              </Button>
            </VStack>
          </Box>
        </Card.Body>
      </Card.Root>
    </VStack>
  );
}
